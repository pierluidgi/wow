package client

import (
	"bufio"
	"errors"
	"github.com/vmihailenco/msgpack/v5"
	"net"
	"time"
	"wow/hashcash"
	"wow/protocol"
)

type Client struct {
	ServerAddr   string
	ConnTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c *Client) request(message *protocol.Message) (*protocol.Message, error) {
	conn, err := net.DialTimeout("tcp", c.ServerAddr, c.ConnTimeout)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	cmdBytes, err := msgpack.Marshal(message)

	if err != nil {
		return nil, err
	}

	if err := conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout)); err != nil {
		return nil, err
	}

	if _, err := conn.Write(cmdBytes); err != nil {
		return nil, err
	}

	if err := conn.SetReadDeadline(time.Now().Add(c.ReadTimeout)); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(conn)
	decoder := msgpack.NewDecoder(reader)

	var response protocol.Message

	if err := decoder.Decode(&response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *Client) requestQuoteWithPoW(data []byte) (string, error) {
	var options protocol.ChallengeOptions

	if err := msgpack.Unmarshal(data, &options); err != nil {
		return "", err
	}

	hc := hashcash.HashCash{
		TargetBits: options.TargetBits,
		Timestamp:  options.Timestamp,
		Data:       options.Data,
		Signature:  options.Signature,
	}

	options.Counter = hc.FindProofCounter()

	data, err := msgpack.Marshal(&options)

	if err != nil {
		return "", err
	}

	return c.requestQuote(&protocol.Message{
		Type: protocol.RequestQuoteWithProof,
		Data: data,
	})
}

func (c *Client) requestQuote(message *protocol.Message) (string, error) {
	resp, err := c.request(message)

	if err != nil {
		return "", err
	}

	switch resp.Type {
	case protocol.ResponseQuote:
		return string(resp.Data), nil
	case protocol.ResponseChallenge:
		if message.Type == protocol.RequestQuote {
			return c.requestQuoteWithPoW(resp.Data)
		}

		return "", errors.New("bad response")
	case protocol.ResponseError:
		return "", errors.New(string(resp.Data))
	}

	return "", errors.New("unknown response")
}

func (c *Client) GetQuote() (string, error) {
	return c.requestQuote(&protocol.Message{
		Type: protocol.RequestQuote,
		Data: nil,
	})
}
