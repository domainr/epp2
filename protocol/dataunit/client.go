package dataunit

import (
	"context"
	"sync"
)

// Client provides an ordered queue of client requests on a data unit connection.
// Each call to ExchangeDataUnit will block until the peer responds, the Context is canceled,
// or the underlying connection is closed. Requests will be processed in strict FIFO order.
// A Client is safe to call from multiple goroutines.
type Client struct {
	writing sync.Mutex
	reading sync.Mutex
	Conn    Conn

	queueing sync.Mutex
	queue    []chan<- result
}

// ExchangeDataUnit sends data unit req and returns the response from the server.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed. The supplied Context must be non-nil.
// Exchange is safe to call from multiple goroutines.
func (c *Client) ExchangeDataUnit(ctx context.Context, req []byte) ([]byte, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	ch := make(chan result, 1)
	go func() {
		err := c.write(ch, req)
		if err != nil {
			ch <- result{nil, err}
			return
		}
		c.read()
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		return res.data, res.err
	}
}

func (c *Client) write(ch chan<- result, data []byte) error {
	c.writing.Lock()
	defer c.writing.Unlock()
	err := c.Conn.WriteDataUnit(data)
	if err != nil {
		return err
	}
	c.enqueue(ch)
	return nil
}

func (c *Client) read() {
	c.reading.Lock()
	defer c.reading.Unlock()
	ch := c.dequeue()
	data, err := c.Conn.ReadDataUnit()
	ch <- result{data, err}
}

func (c *Client) enqueue(ch chan<- result) {
	c.queueing.Lock()
	defer c.queueing.Unlock()
	c.queue = append(c.queue, ch)
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
