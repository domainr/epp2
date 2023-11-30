// Package dataunit implements low-level encoding of the EPP data unit protocol.
//
// The data unit protocol is a simple framing of an XML payload prefixed with a 4-byte header
// in network byte order that expresses the total length of the data unit (header + payload size).
// EPP data units are sent and received via an underlying transport (typically a TLS connection).
//
// See https://datatracker.ietf.org/doc/rfc4934/ for more information.
package dataunit

import (
	"io"
	"net"
)

// Reader is the interface implemented by any type that can read an EPP data unit.
type Reader interface {
	// ReadDataUnit reads a single EPP data unit, returning the payload bytes or an error.
	ReadDataUnit() ([]byte, error)
}

// Writer is the interface implemented by any type that can write an EPP data unit.
type Writer interface {
	// WriteDataUnit writes a single EPP data unit, returning any error.
	WriteDataUnit([]byte) error
}

// Conn is the interface implemented by any type that can read and write EPP data units.
// Concurrent operations on a Conn are implementation-specific and should
// be protected by a synchronization mechanism.
type Conn interface {
	Reader
	Writer
	Close() error
}

// Pipe returns two [Conn] instances that represent the two endpoints of an EPP connection.
// It uses [io.Pipe] to synchronize reads and writes. Each Conn returned is identical,
// either can be used as a client or a server endpoint.
func Pipe() (Conn, Conn) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	return &pipe{r1, w2}, &pipe{r2, w1}
}

// pipe implements [Conn] using an io.Reader and an io.Writer.
type pipe struct {
	R io.Reader
	W io.Writer
}

var _ Conn = &pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (p *pipe) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(p.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (p *pipe) WriteDataUnit(data []byte) error {
	return WriteDataUnit(p.W, data)
}

// Close attempts to close both the underlying reader and writer.
// It will return the first error encountered.
func (p *pipe) Close() error {
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

// ReadDataUnit reads a single EPP data unit from c, returning the payload or an error.
func (c *NetConn) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to c or returns an error.
func (c *NetConn) WriteDataUnit(data []byte) error {
	return WriteDataUnit(c.Conn, data)
}
