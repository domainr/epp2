// Package wire implements low-level encoding of the EPP wire protocol.
//
// The EPP wire protocol is a simple framing of UTF-8 encoded XML prefixed with a 4-byte,
// big-endian header that expresses the total length of the EPP data unit (header + payload size).
// EPP data units are sent and received via an underlying transport (typically a TLS connection).
//
// See https://datatracker.ietf.org/doc/rfc4934/ for more information.
package wire

import (
	"encoding/binary"
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

// ReadDataUnit reads a single EPP data unit from c, returning the payload or an error.
func (c *NetConn) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(c.Conn)
}

// WriteDataUnit writes a single EPP data unit to c or returns an error.
func (c *NetConn) WriteDataUnit(data []byte) error {
	return WriteDataUnit(c.Conn, data)
}

// ReadDataUnit reads a single EPP data unit from [io.Reader] r, returning the payload or an error.
//
// An EPP data unit is prefixed with 32-bit, big-endian value specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See https://datatracker.ietf.org/doc/rfc4934/ for more information.
func ReadDataUnit(r io.Reader) ([]byte, error) {
	var n uint32
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}
	// An EPP data unit size includes the 4 byte header.
	// See https://tools.ietf.org/html/rfc5734#section-4.
	if n < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	n -= 4
	p := make([]byte, n)
	_, err = io.ReadAtLeast(r, p, int(n))
	return p, err
}

// WriteDataUnit writes a single EPP data unit to [io.Writer] w.
//
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See https://datatracker.ietf.org/doc/rfc4934/ for more information.
func WriteDataUnit(w io.Writer, p []byte) error {
	s := uint32(4 + len(p))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(p)
	return err
}
