package dataunit

import (
	"context"
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
	Next(context.Context) ([]byte, Writer, error)
	Close() error
}

type server struct {
	ctx     context.Context
	cancel  func()
	reading sync.Mutex
	writing sync.Mutex
	conn    Conn
	pending chan transaction
}

func NewServer(conn Conn, depth int) Server {
	// TODO: accept a Context
	ctx, cancel := context.WithCancel(context.Background())

	c := &server{
		ctx:     ctx,
		cancel:  cancel,
		conn:    conn,
		pending: make(chan transaction, depth),
	}

	return c
}

func (c *server) Close() error {
	// TODO: gracefully terminate?
	c.cancel()
	return c.conn.Close()
}

func (c *server) Next(ctx context.Context) ([]byte, Writer, error) {
	return c.read(ctx)
}

func (c *server) read(ctx context.Context) ([]byte, Writer, error) {
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
	case <-ctx.Done():
		return nil, f, ctx.Err()
	case <-c.ctx.Done():
		return nil, f, c.ctx.Err()
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
		case <-c.ctx.Done():
			// TODO: do something graceful here?
			return c.ctx.Err()
		case tx = <-c.pending:
		default:
			// Nothing queued, return
			return nil
		}

		var res []byte
		select {
		case <-c.ctx.Done():
			// TODO: do something graceful here?
			err := c.ctx.Err()
			tx.err <- err
			return err
		case res = <-tx.res:
		}

		c.writing.Lock()
		err := c.conn.WriteDataUnit(res)
		c.writing.Unlock()
		if err != nil {
			tx.err <- err
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
