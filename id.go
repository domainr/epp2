package epp

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"sync/atomic"
)

type seqSource struct {
	prefix string
	n      atomic.Uint64
}

func newSeqSource(prefix string) (*seqSource, error) {
	if prefix == "" {
		var pfx [16]byte
		_, err := rand.Read(pfx[:])
		if err != nil {
			return nil, err
		}
		prefix = hex.EncodeToString(pfx[:])
	}
	return &seqSource{
		prefix: prefix,
	}, nil
}

func (s *seqSource) ID() string {
	return s.prefix + strconv.FormatUint(s.id(), 10)
}

func (s *seqSource) id() uint64 {
	return s.n.Add(1)
}
