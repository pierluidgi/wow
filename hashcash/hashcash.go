package hashcash

import (
	"crypto/sha1"
	"fmt"
	"math/big"
)

type HashCash struct {
	TargetBits byte
	Timestamp  uint32
	Data       uint64
	Signature  uint32
	Counter    uint32
}

func (hc HashCash) String() string {
	return fmt.Sprintf("%d:%d:%d:%d:%d", hc.TargetBits, hc.Timestamp, hc.Data, hc.Signature, hc.Counter)
}

func (hc HashCash) ZeroBits() byte {
	digest := sha1.Sum([]byte(hc.String()))
	digestHex := new(big.Int).SetBytes(digest[:])
	return byte((sha1.Size * 8) - digestHex.BitLen())
}

func (hc *HashCash) FindProofCounter() uint32 {
	hc.Counter = 0

	for hc.ZeroBits() < hc.TargetBits {
		hc.Counter++
	}

	return hc.Counter
}
