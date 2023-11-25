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
	Close() error
}

type server struct {
	reading sync.Mutex
	writing sync.Mutex
	conn    Conn
	pending chan transaction
	closed  chan struct{}
}

func NewServer(conn Conn, depth int) Server {
	return &server{
		conn:    conn,
		pending: make(chan transaction, depth),
		closed:  make(chan struct{}),
	}
}

func (c *server) Close() error {
	// TODO: gracefully terminate?
	close(c.closed)
	return c.conn.Close()
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
		err := c.writePending()
		if err != nil {
			return err
		}
		return <-tx.err
	})
	// TODO: enqueue transaction before or after reading from conn?
	select {
	case <-c.closed:
		return nil, f, ErrClosedConnection
	case c.pending <- tx:
	}
	c.reading.Lock()
	data, err := c.conn.ReadDataUnit()
	c.reading.Unlock()
	return data, f, err
}

func (c *server) writePending() error {
	for {
		var tx transaction
		select {
		case <-c.closed:
			return ErrClosedConnection
		case tx = <-c.pending:
		default:
			// Nothing queued, return
			return nil
		}

		var res []byte
		select {
		case <-c.closed:
			return ErrClosedConnection
		case res = <-tx.res:
		}

		c.writing.Lock()
		err := c.conn.WriteDataUnit(res)
		c.writing.Unlock()
		tx.err <- err
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
