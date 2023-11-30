package dataunit

import (
	"sync"
)

// Client provides an ordered queue of client requests on a data unit connection.
// Each call to SendDataUnit will block until the peer responds or the underlying connection is closed.
// Requests will be processed in strict FIFO order.
// A Client is safe to call from multiple goroutines.
type Client interface {
	SendDataUnit([]byte) ([]byte, error)
}

type client struct {
	writing sync.Mutex
	pending []chan result

	reading sync.Mutex

	// Reads and writes are protected by reading and writing, respectively.
	conn Conn
}

func NewClient(conn Conn) Client {
	return &client{conn: conn}
}

func (c *client) SendDataUnit(data []byte) ([]byte, error) {
	return c.exchange(data)
}

func (c *client) exchange(data []byte) ([]byte, error) {
	head, tail, err := c.write(data)
	if err != nil {
		return nil, err
	}
	c.read(head)
	res := <-tail
	return res.data, res.err
}

func (c *client) write(data []byte) (head chan<- result, tail <-chan result, err error) {
	c.writing.Lock()
	defer c.writing.Unlock()
	err = c.conn.WriteDataUnit(data)
	if err != nil {
		return nil, nil, err
	}

	// TODO: optimize this
	chtail := make(chan result, 1)
	c.pending = append(c.pending, chtail)
	head = c.pending[0]
	c.pending = c.pending[1:]

	return head, chtail, err
}

func (c *client) read(ch chan<- result) {
	c.reading.Lock()
	defer c.reading.Unlock()
	data, err := c.conn.ReadDataUnit()
	ch <- result{data, err}
}

type result struct {
	data []byte
	err  error
}
