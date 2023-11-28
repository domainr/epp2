package dataunit

import (
	"sync"
)

// Receiver provides an ordered queue of client requests coupled with a [Writer]
// to respond to the request. The Writer returned from ReceiveDataUnit can be called once.
// Calling WriteDataUnit more than once is undefined.
// Receiver enforces ordering of responses, writing each response in the same
// order as the requests received from the client.
// A Receiver is safe to call from multiple goroutines, and each client request
// may be handled in a separate goroutine.
type Receiver interface {
	ReceiveDataUnit() ([]byte, Writer, error)
}

type receiver struct {
	reading sync.Mutex
	reads   uint64

	writing sync.Mutex
	writes  uint64
	pending []transaction

	// Reads and writes are protected by reading and writing, respectively.
	conn Conn
}

func NewReceiver(conn Conn) Receiver {
	return &receiver{conn: conn}
}

func (r *receiver) ReceiveDataUnit() ([]byte, Writer, error) {
	return r.read()
}

func (r *receiver) read() ([]byte, Writer, error) {
	r.reading.Lock()
	defer r.reading.Unlock()

	n := r.reads
	r.reads += 1

	f := writerFunc(func(data []byte) error {
		ch, err := r.respond(n, data)
		if ch == nil {
			return err
		}
		return <-ch
	})
	data, err := r.conn.ReadDataUnit()
	return data, f, err
}

func (r *receiver) respond(n uint64, data []byte) (<-chan error, error) {
	r.writing.Lock()
	defer r.writing.Unlock()

	i := int(n - r.writes)

	// If this isn’t the oldest pending transaction, queue the response.
	if i > 0 {
		if i > len(r.pending) {
			r.pending = append(r.pending, make([]transaction, i-len(r.pending))...)
		}
		ch := make(chan error, 1)
		r.pending[i-1] = transaction{data, ch}
		return ch, nil
	}

	// Write responses
	err := r.conn.WriteDataUnit(data)
	if err != nil {
		return nil, err
	}
	r.writes += 1
	var writes int
	for _, tx := range r.pending {
		if tx.res == nil {
			break
		}
		err := r.conn.WriteDataUnit(data)
		tx.err <- err
		if err != nil {
			break
		}
		writes += 1
	}
	if writes == len(r.pending) {
		r.pending = r.pending[:0:min(cap(r.pending), capMax)]
	} else {
		r.pending = r.pending[writes:]
	}
	r.writes += uint64(writes)

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
