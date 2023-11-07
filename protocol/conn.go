package protocol

import (
	"io"
	"net"
)

// Conn is the interface implemented by any type that can read and write EPP data units.
// Concurrent operations on a Conn are implementation-specific and should
// be protected by a synchronization mechanism.
type Conn interface {
	// ReadDataUnit reads a single EPP data unit, returning the payload bytes or an error.
	ReadDataUnit() ([]byte, error)

	// WriteDataUnit writes a single EPP data unit, returning any error.
	WriteDataUnit([]byte) error

	// Close closes the connection.
	Close() error
}

// Pipe implements [Conn] using an io.Reader and an io.Writer.
type Pipe struct {
	R io.Reader
	W io.Writer
}

var _ Conn = &Pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (p *Pipe) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(p.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (p *Pipe) WriteDataUnit(data []byte) error {
	return WriteDataUnit(p.W, data)
}

// Close attempts to close both the underlying reader and writer.
// It will return the first error encountered.
func (p *Pipe) Close() error {
	var rerr, werr error
	if c, ok := p.R.(io.Closer); ok {
		rerr = c.Close()
	}
	if r, ok := p.W.(io.Reader); ok && r == p.R {
		return rerr
	}
	if c, ok := p.W.(io.Closer); ok {
		werr = c.Close()
	}
	if rerr != nil {
		return rerr
	}
	return werr
}

// NetConn implements [Conn] using a net.Conn.
type NetConn struct {
	net.Conn
}

var _ Conn = &NetConn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (c *NetConn) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (c *NetConn) WriteDataUnit(data []byte) error {
	return WriteDataUnit(c.Conn, data)
}
