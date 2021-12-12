package transport

import (
	"io"
	"net"
	"sync"
)

// Transport is a generic connection that can read and write EPP data units.
// Multiple goroutines may invoke methods on a Transport simultaneously.
type Transport interface {
	// ReadDataUnit reads a single EPP data unit, returning the payload bytes or an error.
	ReadDataUnit() ([]byte, error)

	// WriteDataUnit writes a single EPP data unit, returning any error.
	WriteDataUnit([]byte) error

	// Close closes the connection.
	Close() error
}

// Pipe implements Transport using an io.Reader and an io.Writer.
type Pipe struct {
	// R is from by ReadDataUnit.
	R io.Reader

	// W is written to by WriteDataUnit.
	W io.Writer

	r sync.Mutex
	w sync.Mutex
}

var _ Transport = &Pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (t *Pipe) ReadDataUnit() ([]byte, error) {
	t.r.Lock()
	defer t.r.Unlock()
	return ReadDataUnit(t.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *Pipe) WriteDataUnit(data []byte) error {
	t.w.Lock()
	defer t.w.Unlock()
	return WriteDataUnit(t.W, data)
}

// Close attempts to close both the underlying reader and writer.
// It will return the first error encountered.
func (t *Pipe) Close() error {
	var rerr, werr error
	if c, ok := t.R.(io.Closer); ok {
		rerr = c.Close()
	}
	if r, ok := t.W.(io.Reader); ok && r == t.R {
		return rerr
	}
	if c, ok := t.W.(io.Closer); ok {
		werr = c.Close()
	}
	if rerr != nil {
		return rerr
	}
	return werr
}

// Conn implements Transport using a net.Conn.
type Conn struct {
	net.Conn
	r sync.Mutex
	w sync.Mutex
}

var _ Transport = &Conn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (c *Conn) ReadDataUnit() ([]byte, error) {
	c.r.Lock()
	defer c.r.Unlock()
	return ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (c *Conn) WriteDataUnit(p []byte) error {
	c.w.Lock()
	defer c.w.Unlock()
	return WriteDataUnit(c.Conn, p)
}
