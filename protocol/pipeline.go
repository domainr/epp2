package protocol

import (
	"context"
	"encoding/xml"
	"sync"

	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

type Option interface{}

type Writer interface {
	WriteEPP(context.Context, epp.Body) error
}

type Request interface {
	Context() context.Context // the underlying server has a base context for graceful shutdown
	Body() epp.Body
	Respond(epp.Body) error
}

type ServerConn interface {
	Next() (Request, error)
	Close(context.Context) error
}

type serverConn struct {
	ctx     context.Context
	cancel  func()
	conn    dataunit.Conn
	schemas schema.Schemas
	idle    chan request
	pending chan request
}

func NewServer(conn dataunit.Conn, opts ...Option) ServerConn {
	const depth = 10 // depth must be >= 1

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

func (c *serverConn) Close(ctx context.Context) error {
	// TODO: gracefully terminate?
	c.cancel()
	return c.conn.Close()
}

func (c *serverConn) Next() (Request, error) {
	// Read EPP message
	// Pack into a Request
	// ...
	// on Request.Respond():
	// wait for turn in queue
	// write response EPP message
	// backpressure here?
	return nil, nil
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
			req.res = make(chan epp.Body, 1)
		}
		req.ctx = c.ctx

		data, err := c.conn.ReadDataUnit()
		if err != nil {
			req.req = func() (epp.Body, error) {
				return nil, err
			}
		} else {
			var once sync.Once
			var body epp.Body
			var err error
			req.req = func() (epp.Body, error) {
				once.Do(func() {
					var e epp.EPP
					err = schema.Unmarshal(data, &e, c.schemas)
					body = e.Body
				})
				return body, err
			}
		}

		// Enqueue the request to write the response.
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

		var e epp.EPP
		select {
		case <-req.ctx.Done():
			// TODO: do something graceful here?
			return
		case e.Body = <-req.res:
		}

		data, err := xml.Marshal(&e)
		if err != nil {
			// TODO: handle error marshaling
			// FIXME: this should return something to the client, even if just "I give up"
			return
		}

		err = c.conn.WriteDataUnit(data)
		if err != nil {
			// TODO: handle error gracefully
			return
		}
	}
}

type request struct {
	ctx context.Context
	req func() (epp.Body, error)
	res chan epp.Body
}
