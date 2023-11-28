package dataunit

import (
	"sync"
)

// Server provides an ordered queue of client requests coupled with a [Writer]
// to respond to the request. The Writer returned from Next can be called once.
// Calling WriteDataUnit more than once is undefined.
// Server enforces ordering of responses, writing each response in the same
// order as the requests received from the client.
// A Server is safe to call from multiple goroutines. Each client request
// may be handled in a separate goroutine.
type Server interface {
	Next() ([]byte, Writer, error)
}

type server struct {
	reading sync.Mutex
	writing sync.Mutex
	conn    Conn

	mu      sync.Mutex
	in      uint64
	out     uint64
	pending []transaction
}

func NewServer(conn Conn) Server {
	return &server{conn: conn}
}

func (c *server) Next() ([]byte, Writer, error) {
	return c.read()
}

func (c *server) read() ([]byte, Writer, error) {
	c.reading.Lock()
	defer c.reading.Unlock()

	c.mu.Lock()
	n := c.in
	c.in += 1
	c.mu.Unlock()

	f := writerFunc(func(data []byte) error {
		ch, err := c.respond(n, data)
		if ch == nil {
			return err
		}
		return <-ch
	})
	data, err := c.conn.ReadDataUnit()
	return data, f, err
}

func (c *server) respond(n uint64, data []byte) (<-chan error, error) {
	const capMax = 32

	c.mu.Lock()
	defer c.mu.Unlock()
	depth := int(c.in - c.out - 1)
	n = n - c.out
	i := int(n - c.out)

	// If this isn’t the oldest pending transaction, queue the response.
	if i > 0 {
		if depth > len(c.pending) {
			c.pending = append(c.pending, make([]transaction, depth)...)
		}
		ch := make(chan error, 1)
		c.pending[i-1] = transaction{data, ch}
		return ch, nil
	}

	// Write responses
	c.writing.Lock()
	defer c.writing.Unlock()
	err := c.conn.WriteDataUnit(data)
	if err != nil {
		return nil, err
	}
	c.out += 1
	var writes uint64
	for _, tx := range c.pending {
		if tx.res == nil {
			break
		}
		err := c.conn.WriteDataUnit(data)
		tx.err <- err
		if err != nil {
			break
		}
		writes += 1
	}
	c.pending = c.pending[writes:len(c.pending):min(cap(c.pending), capMax)]
	c.out += writes

	return nil, nil
}

type transaction struct {
	res []byte
	err chan error
}

type writerFunc func(data []byte) error

func (f writerFunc) WriteDataUnit(data []byte) error {
	return f(data)
}
