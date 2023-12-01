package dataunit

import (
	"sync"
)

// Client provides an ordered queue of client requests on a data unit connection.
// Each call to ExchangeDataUnit will block until the peer responds or the underlying connection is closed.
// Requests will be processed in strict FIFO order.
// A Client is safe to call from multiple goroutines.
type Client struct {
	writing sync.Mutex
	reading sync.Mutex
	Conn    Conn

	queueing sync.Mutex
	queue    []chan result
}

func (c *Client) ExchangeDataUnit(data []byte) ([]byte, error) {
	ch, err := c.write(data)
	if err != nil {
		return nil, err
	}
	c.read()
	res := <-ch
	return res.data, res.err
}

func (c *Client) write(data []byte) (<-chan result, error) {
	c.writing.Lock()
	defer c.writing.Unlock()
	ch := c.enqueue()
	err := c.Conn.WriteDataUnit(data)
	return ch, err
}

func (c *Client) read() {
	c.reading.Lock()
	defer c.reading.Unlock()
	ch := c.dequeue()
	data, err := c.Conn.ReadDataUnit()
	ch <- result{data, err}
}

func (c *Client) enqueue() <-chan result {
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
