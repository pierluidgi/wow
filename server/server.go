package server

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"math/rand"
	"net"
	"sync"
	"time"
	"wow/cache"
	"wow/hashcash"
	"wow/metrics"
	"wow/protocol"
	"wow/storage"
)

type Options struct {
	Address       string
	ReadTimeout   int
	WriteTimeout  int
	DDoSRate      uint64
	TargetBits    byte
	ChallengeTtl  uint32
	QuotesStorage storage.QuotesStorage
	Cache         cache.Cache
	RateMeter     *metrics.RateMeter
}

type Server struct {
	address       string
	readTimeout   time.Duration
	writeTimeout  time.Duration
	ddosRate      uint64
	targetBits    byte
	challengeTtl  uint32
	quotesStorage storage.QuotesStorage
	cache         cache.Cache
	rateMeter     *metrics.RateMeter
	secretKey     uint32
	msgPool       sync.Pool
}

func (s *Server) getMessage() *protocol.Message {
	v := s.msgPool.Get()

	if v == nil {
		return &protocol.Message{
			Type: 0,
			Data: nil,
		}
	}

	return v.(*protocol.Message)
}

func (s *Server) putMessage(message *protocol.Message) {
	message.Data = message.Data[:0]
	s.msgPool.Put(message)
}

func NewServer(options *Options) *Server {
	rand.Seed(time.Now().UnixNano())

	return &Server{
		address:       options.Address,
		readTimeout:   time.Duration(options.ReadTimeout) * time.Second,
		writeTimeout:  time.Duration(options.WriteTimeout) * time.Second,
		ddosRate:      options.DDoSRate,
		targetBits:    options.TargetBits,
		challengeTtl:  options.ChallengeTtl,
		quotesStorage: options.QuotesStorage,
		cache:         options.Cache,
		rateMeter:     options.RateMeter,
		secretKey:     rand.Uint32(),
		msgPool:       sync.Pool{},
	}
}

func (s *Server) writeResponse(conn net.Conn, message *protocol.Message) {
	respBytes, err := msgpack.Marshal(message)

	if err != nil {
		log.Error(err)
		return
	}

	if err := conn.SetWriteDeadline(time.Now().Add(s.writeTimeout)); err != nil {
		log.Error(err)
		return
	}

	if _, err := conn.Write(respBytes); err != nil {
		log.Error(respBytes)
	}
}

func (s *Server) respondError(conn net.Conn, message string) {
	s.writeResponse(conn, &protocol.Message{
		Type: protocol.ResponseError,
		Data: []byte(message),
	})
}

func (s *Server) respondQuote(conn net.Conn) {
	quote, err := s.quotesStorage.RandomQuote()

	if err != nil {
		log.Error(err)
		s.respondError(conn, protocol.InternalServerErrorMsg)
		return
	}

	s.writeResponse(conn, &protocol.Message{
		Type: protocol.ResponseQuote,
		Data: []byte(quote),
	})
}

func (s *Server) challengeSignature(hc *hashcash.HashCash) uint32 {
	signature := uint32(hc.TargetBits) ^ hc.Timestamp ^ uint32(hc.Data>>32) ^ uint32(hc.Data>>16) ^ s.secretKey
	return signature
}

func (s *Server) generateChallenge() *protocol.ChallengeOptions {
	hc := hashcash.HashCash{
		TargetBits: s.targetBits,
		Timestamp:  uint32(time.Now().Unix()),
		Data:       rand.Uint64(),
		Signature:  0,
		Counter:    0,
	}

	hc.Signature = s.challengeSignature(&hc)

	return &protocol.ChallengeOptions{
		TargetBits: hc.TargetBits,
		Timestamp:  hc.Timestamp,
		Data:       hc.Data,
		Signature:  hc.Signature,
		Counter:    hc.Counter,
	}
}

func (s *Server) respondChallenge(conn net.Conn) {
	data := s.generateChallenge()

	dataBytes, err := msgpack.Marshal(&data)

	if err != nil {
		log.Error(err)
		return
	}

	s.writeResponse(conn, &protocol.Message{
		Type: protocol.ResponseChallenge,
		Data: dataBytes,
	})
}

func (s *Server) validateProof(conn net.Conn, data []byte) {
	var options protocol.ChallengeOptions

	if err := msgpack.Unmarshal(data, &options); err != nil {
		s.respondError(conn, protocol.BadRequestErrorMsg)
		return
	}

	currTime := uint32(time.Now().Unix())

	if options.Timestamp+s.challengeTtl < currTime {
		s.respondError(conn, protocol.ChallengeExpiredErrorMsg)
		return
	}

	hc := hashcash.HashCash{
		TargetBits: options.TargetBits,
		Timestamp:  options.Timestamp,
		Data:       options.Data,
		Signature:  options.Signature,
		Counter:    options.Counter,
	}

	if options.Signature != s.challengeSignature(&hc) {
		s.respondError(conn, protocol.BadProofErrorMsg)
		return
	}

	if hc.ZeroBits() < options.TargetBits {
		s.respondError(conn, protocol.BadProofErrorMsg)
		return
	}

	if contains := s.cache.ContainsOrAdd(options.Data ^ (uint64(options.Timestamp) << 32)); contains {
		s.respondError(conn, protocol.BadProofErrorMsg)
		return
	}

	s.respondQuote(conn)
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	s.rateMeter.Add(1)

	if err := conn.SetReadDeadline(time.Now().Add(s.readTimeout)); err != nil {
		log.Error(err)
		return
	}

	reader := bufio.NewReader(conn)
	decoder := msgpack.NewDecoder(reader)

	msg := s.getMessage()
	defer s.putMessage(msg)

	if err := decoder.Decode(&msg); err != nil {
		log.Error(err)
		return
	}

	switch msg.Type {
	case protocol.RequestQuote:
		if s.rateMeter.Rate() < s.ddosRate {
			s.respondQuote(conn)
			return
		}
		s.respondChallenge(conn)
	case protocol.RequestQuoteWithProof:
		s.validateProof(conn, msg.Data)
	default:
		s.respondError(conn, protocol.BadRequestErrorMsg)
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.address)

	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			log.Error(err)
			continue
		}

		go s.handleConn(conn)
	}
}
