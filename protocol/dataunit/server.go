package dataunit

import (
	"context"
)

type ServerConn interface {
	Next(context.Context) ([]byte, Writer, error)
	Close() error
}

type serverConn struct {
	ctx     context.Context
	cancel  func()
	conn    Conn
	idle    chan request
	pending chan request
}

func NewServer(conn Conn, depth int) ServerConn {
	ctx, cancel := context.WithCancel(context.Background())

	c := &serverConn{
		ctx:     ctx,
		cancel:  cancel,
		conn:    conn,
		idle:    make(chan request, depth),
		pending: make(chan request, depth),
	}

	for i := 0; i < depth; i++ {
		c.idle <- request{}
	}

	go c.read()
	go c.write()
	return c
}

func (c *serverConn) Close() error {
	// TODO: gracefully terminate?
	c.cancel()
	return c.conn.Close()
}

func (c *serverConn) Next(ctx context.Context) ([]byte, Writer, error) {
	// Read EPP message
	// Pack into a Request
	// ...
	// on Request.Respond():
	// wait for turn in queue
	// write response EPP message
	// backpressure here?
	return nil, nil, nil
}

func (c *serverConn) work() {
	var current request
	for {
		select {
		case <-c.ctx.Done():
			return
		case res := <-current.res:
			// Write response
			err := c.conn.WriteDataUnit(res)
			if err != nil {
				// TODO: handle error gracefully
				return
			}
			c.idle <- current
			select {
			case current = <-c.pending:
			}
		case req := <-c.idle:
			if req.res == nil {
				req.res = make(chan []byte, 1)
			}
			req.ctx = c.ctx
			req.data, req.err = c.conn.ReadDataUnit()
			if current.res == nil {
				current = req
			} else {
				c.pending <- req
			}
		}
	}
}

func (c *serverConn) read() {
	for {
		var req request
		select {
		case <-c.ctx.Done():
			// TODO: respond to any pending requests with a shutting down error?
			return
		case req = <-c.idle:
		}

		if req.res == nil {
			req.res = make(chan []byte, 1)
		}
		req.ctx = c.ctx
		req.data, req.err = c.conn.ReadDataUnit()

		// Enqueue req for writing a response.
		// This should never block.
		c.pending <- req
	}
}

func (c *serverConn) write() {
	for {
		var req request
		select {
		case <-c.ctx.Done():
			// TODO: do something graceful here?
			return
		case req = <-c.pending:
		}

		var res []byte
		select {
		case <-req.ctx.Done():
			// TODO: do something graceful here?
			return
		case res = <-req.res:
		}

		err := c.conn.WriteDataUnit(res)
		if err != nil {
			// TODO: handle error gracefully
			return
		}

		// Return req to the idle queue.
		// This should never block.
		c.idle <- req
	}
}

type request struct {
	ctx  context.Context
	data []byte
	err  error
	res  chan []byte
}

func (r *request) WriteDataUnit(data []byte) error {
	r.res <- data
	return nil
}
