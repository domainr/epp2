package dataunit

import (
	"context"
	"io"
	"strconv"
	"sync"
	"sync/atomic"
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

// MultipleResponseError is returned when a [Responder] is called more than once.
type MultipleResponseError struct {
	Index uint64
	Count uint64
}

func (err MultipleResponseError) Error() string {
	return "epp: multiple responses to request " + strconv.FormatUint(err.Index, 10) +
		": " + strconv.FormatUint(err.Count, 10) + " > 1"
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

	Conn io.ReadWriter
}

// ServeDataUnit reads one data unit from the client and provides a [Responder] to respond.
//
// The supplied Context must be non-nil, and only affects reading the request from the client.
// Cancelling the Context after ServeDataUnit returns will have no effect on the Responder.
//
// The returned Responder can only be called once. The returned Responder will always
// be non-nil, so the caller can respond to a malformed client request.
//
// ServeDataUnit is safe to be called from multiple goroutines, and each client request
// may be handled in a separate goroutine.
func (s *Server) ServeDataUnit(ctx context.Context) ([]byte, Responder, error) {
	s.reading.Lock()
	defer s.reading.Unlock()

	n := s.reads
	s.reads += 1

	var counter atomic.Uint64

	f := responderFunc(func(ctx context.Context, data []byte) error {
		count := counter.Add(1)
		if count != 1 {
			return MultipleResponseError{Index: n, Count: count}
		}
		ch, err := s.respond(ctx, n, data)
		if ch == nil {
			return err
		}
		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case err = <-ch:
			return err
		}
	})

	data, err := Receive(ctx, s.Conn)
	return data, f, err
}

func (s *Server) respond(ctx context.Context, n uint64, data []byte) (<-chan error, error) {
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
	err := Send(ctx, s.Conn, data)
	if err != nil {
		return nil, err
	}
	s.writes += 1
	var writes int
	for _, tx := range s.pending {
		if tx.res == nil {
			break
		}
		err := Send(ctx, s.Conn, tx.res)
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
