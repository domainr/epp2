package dataunit

import (
	"fmt"
	"sync"
)

// Client provides an ordered queue of client requests on a data unit connection.
// Each call to SendDataUnit will block until the peer responds or the underlying connection is closed.
// A Client is safe to call from multiple goroutines. Requests will be processed in strict FIFO order.
type Client interface {
	SendDataUnit([]byte) ([]byte, error)
}

type client struct {
	writing sync.Mutex
	reading sync.Mutex
	conn    Conn

	queueing sync.Mutex
	queue    []chan result
}

func NewClient(conn Conn) Client {
	return &client{conn: conn}
}

func (c *client) SendDataUnit(data []byte) ([]byte, error) {
	return c.exchange(data)
}

func (c *client) exchange(data []byte) ([]byte, error) {
	ch, err := c.write(data)
	if err != nil {
		return nil, err
	}
	c.read()
	res := <-ch
	fmt.Printf("client: %s %s\n", string(data), string(res.data))
	return res.data, res.err
}

func (c *client) write(data []byte) (<-chan result, error) {
	c.writing.Lock()
	defer c.writing.Unlock()
	err := c.conn.WriteDataUnit(data)
	if err != nil {
		return nil, err
	}
	return c.enqueue(), nil
}

func (c *client) read() {
	c.reading.Lock()
	defer c.reading.Unlock()
	ch := c.dequeue()
	data, err := c.conn.ReadDataUnit()
	ch <- result{data, err}
}

func (c *client) enqueue() <-chan result {
	c.queueing.Lock()
	defer c.queueing.Unlock()
	ch := make(chan result, 1)
	c.queue = append(c.queue, ch)
	return ch
}

func (c *client) dequeue() chan<- result {
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
