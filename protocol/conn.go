package protocol

import (
	"io"
	"net"

	"github.com/domainr/epp2/protocol/wire"
)

// Pipe implements [wire.Interface] using an io.Reader and an io.Writer.
type Pipe struct {
	R io.Reader
	W io.Writer
}

var _ wire.Conn = &Pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (p *Pipe) ReadDataUnit() ([]byte, error) {
	return wire.ReadDataUnit(p.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (p *Pipe) WriteDataUnit(data []byte) error {
	return wire.WriteDataUnit(p.W, data)
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

// Conn implements [wire.Interface] using a net.Conn.
type Conn struct {
	net.Conn
}

var _ wire.Conn = &Conn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (c *Conn) ReadDataUnit() ([]byte, error) {
	return wire.ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (c *Conn) WriteDataUnit(data []byte) error {
	return wire.WriteDataUnit(c.Conn, data)
}
