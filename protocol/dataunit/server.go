package dataunit

import (
	"context"
	"sync"
)

// Responder is the interface implemented by any type that can respond to a data unit request.
type Responder interface {
	// RespondDataUnit will attempt to write data to the underlying connection.
	// It blocks until the data unit is successfully written, Context is canceled,
	// or the underlying connection is closed. The supplied Context must be non-nil.
	RespondDataUnit(context.Context, []byte) error
}

type responderFunc func(context.Context, []byte) error

func (f responderFunc) RespondDataUnit(ctx context.Context, data []byte) error {
	return f(ctx, data)
}

// Server provides an ordered queue of client requests coupled with a [Responder]
// to respond to the request.
// Server enforces ordering of responses, writing each response in the same
// order as the requests received from the client.
// ServeDataUnit is safe to be called from multiple goroutines.
type Server struct {
	// reading protects Conn and reading
	reading sync.Mutex
	reads   uint64

	// writing protects Conn, writes, and pending
	writing sync.Mutex
	writes  uint64
	pending []transaction

	Conn Conn
}

// ServeDataUnit reads one data unit from the client and provides a [Responder] to respond.
// The returned Writer can only be called once. The returned Responder will always
// be non-nil, so the caller can respond to a malformed client request.
// ServeDataUnit is safe to be called from multiple goroutines, and each client request
// may be handled in a separate goroutine.
func (s *Server) ServeDataUnit(ctx context.Context) ([]byte, Responder, error) {
	s.reading.Lock()
	defer s.reading.Unlock()

	n := s.reads
	s.reads += 1

	f := responderFunc(func(ctx context.Context, data []byte) error {
		ch, err := s.respond(n, data)
		if ch == nil {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-ch:
			return err
		}
	})

	ch := make(chan result, 1)
	go func() {
		data, err := s.Conn.ReadDataUnit()
		ch <- result{data, err}
	}()

	select {
	case <-ctx.Done():
		return nil, f, ctx.Err()
	case res := <-ch:
		return res.data, f, res.err
	}
}

func (s *Server) respond(n uint64, data []byte) (<-chan error, error) {
	s.writing.Lock()
	defer s.writing.Unlock()

	i := int(n - s.writes)

	// If this isn’t the oldest pending transaction, queue the response.
	if i > 0 {
		if i > len(s.pending) {
			s.pending = append(s.pending, make([]transaction, i-len(s.pending))...)
		}
		ch := make(chan error, 1)
		s.pending[i-1] = transaction{data, ch}
		return ch, nil
	}

	// Write responses
	err := s.Conn.WriteDataUnit(data)
	if err != nil {
		return nil, err
	}
	s.writes += 1
	var writes int
	for _, tx := range s.pending {
		if tx.res == nil {
			break
		}
		err := s.Conn.WriteDataUnit(tx.res)
		tx.err <- err
		if err != nil {
			break
		}
		writes += 1
	}
	if writes == len(s.pending) {
		s.pending = s.pending[:0:min(cap(s.pending), capMax)]
	} else {
		s.pending = s.pending[writes:]
	}
	s.writes += uint64(writes)

	return nil, nil
}

const capMax = 32

type transaction struct {
	res []byte
	err chan error
}
