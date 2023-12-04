package dataunit

import (
	"context"
	"io"
	"sync"
)

// Client provides an ordered queue of client requests on a data unit connection.
// Each call to ExchangeDataUnit will block until the peer responds, the Context is canceled,
// or the underlying connection is closed. Requests will be processed in strict FIFO order.
// A Client is safe to call from multiple goroutines.
type Client struct {
	writing sync.Mutex
	reading sync.Mutex
	Conn    io.ReadWriter

	queueing sync.Mutex
	queue    []chan<- result
}

// ExchangeDataUnit sends data unit req and returns the response from the server.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed. The supplied Context must be non-nil.
// Exchange is safe to call from multiple goroutines.
func (c *Client) ExchangeDataUnit(ctx context.Context, data []byte) ([]byte, error) {
	ch, err := c.send(ctx, data)
	if err != nil {
		return nil, err
	}
	go c.receive()
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case res := <-ch:
		return res.data, res.err
	}
}

func (c *Client) send(ctx context.Context, data []byte) (<-chan result, error) {
	c.writing.Lock()
	defer c.writing.Unlock()
	err := Send(ctx, c.Conn, data)
	if err != nil {
		return nil, err
	}
	return c.enqueue(), nil
}

func (c *Client) receive() {
	c.reading.Lock()
	defer c.reading.Unlock()
	ch := c.dequeue()
	data, err := Read(c.Conn)
	ch <- result{data, err}
}

func (c *Client) enqueue() chan result {
	c.queueing.Lock()
	defer c.queueing.Unlock()
	ch := make(chan result, 1)
	c.queue = append(c.queue, ch)
	return ch
}

func (c *Client) dequeue() chan<- result {
	c.queueing.Lock()
	defer c.queueing.Unlock()
	ch := c.queue[0]
	c.queue = c.queue[1:]
	return ch
}

type result struct {
	data []byte
	err  error
}
