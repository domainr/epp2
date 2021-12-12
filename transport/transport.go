package transport

import (
	"io"
	"net"
)

// Transport is a generic connection that can read and write EPP data units.
// Concurrent operations on a Transport are implementation-specific and should
// be protected by a synchronization mechanism.
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
	R io.Reader
	W io.Writer
}

var _ Transport = &Pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (t *Pipe) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(t.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *Pipe) WriteDataUnit(data []byte) error {
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
}

var _ Transport = &Conn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (c *Conn) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (c *Conn) WriteDataUnit(p []byte) error {
	return WriteDataUnit(c.Conn, p)
}
