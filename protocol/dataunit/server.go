package dataunit

import (
	"sync"
)

// Server provides an ordered queue of client requests coupled with a [Writer]
// to respond to the request. The Writer returned from ReceiveDataUnit can be called once.
// Calling WriteDataUnit more than once is undefined.
// Server enforces ordering of responses, writing each response in the same
// order as the requests received from the client.
// A Server is safe to call from multiple goroutines, and each client request
// may be handled in a separate goroutine.
type Server interface {
	ReceiveDataUnit() ([]byte, Writer, error)
}

type server struct {
	reading sync.Mutex
	reads   uint64

	writing sync.Mutex
	writes  uint64
	pending []transaction

	// Reads and writes are protected by reading and writing, respectively.
	conn Conn
}

func NewServer(conn Conn) Server {
	return &server{conn: conn}
}

func (s *server) ReceiveDataUnit() ([]byte, Writer, error) {
	return s.read()
}

func (s *server) read() ([]byte, Writer, error) {
	s.reading.Lock()
	defer s.reading.Unlock()

	n := s.reads
	s.reads += 1

	f := writerFunc(func(data []byte) error {
		ch, err := s.respond(n, data)
		if ch == nil {
			return err
		}
		return <-ch
	})
	data, err := s.conn.ReadDataUnit()
	return data, f, err
}

func (s *server) respond(n uint64, data []byte) (<-chan error, error) {
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
	err := s.conn.WriteDataUnit(data)
	if err != nil {
		return nil, err
	}
	s.writes += 1
	var writes int
	for _, tx := range s.pending {
		if tx.res == nil {
			break
		}
		err := s.conn.WriteDataUnit(data)
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

type writerFunc func(data []byte) error

func (f writerFunc) WriteDataUnit(data []byte) error {
	return f(data)
}
