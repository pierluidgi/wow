package protocol

type MessageType int

const (
	RequestQuote MessageType = iota
	RequestQuoteWithProof
	ResponseQuote
	ResponseChallenge
	ResponseError

	InternalServerErrorMsg   = "internal server error"
	BadRequestErrorMsg       = "bad request"
	BadProofErrorMsg         = "bad proof"
	ChallengeExpiredErrorMsg = "challenge expired"
)

type Message struct {
	Type MessageType
	Data []byte
}

type ChallengeOptions struct {
	TargetBits byte
	Timestamp  uint32
	Data       uint64
	Signature  uint32
	Counter    uint32
}
