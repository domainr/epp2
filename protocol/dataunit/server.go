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
	pending chan transaction
}

func NewServer(conn Conn, depth int) Server {
	return &server{
		conn:    conn,
		pending: make(chan transaction, depth),
	}
}

func (c *server) Next() ([]byte, Writer, error) {
	return c.read()
}

func (c *server) read() ([]byte, Writer, error) {
	tx := transaction{
		res: make(chan []byte, 1),
		err: make(chan error, 1),
	}
	f := writerFunc(func(data []byte) error {
		tx.res <- data
		c.writePending()
		return <-tx.err
	})
	c.pending <- tx // blocks if pipeline is full
	c.reading.Lock()
	data, err := c.conn.ReadDataUnit()
	c.reading.Unlock()
	return data, f, err
}

func (c *server) writePending() {
	for {
		select {
		case tx := <-c.pending:
			select {
			case res := <-tx.res:
				c.writing.Lock()
				err := c.conn.WriteDataUnit(res)
				c.writing.Unlock()
				tx.err <- err
			default:
				return
			}
		default:
			return // nothing queued
		}

	}
}

type transaction struct {
	res chan []byte
	err chan error
}

type writerFunc func(data []byte) error

func (f writerFunc) WriteDataUnit(data []byte) error {
	return f(data)
}
