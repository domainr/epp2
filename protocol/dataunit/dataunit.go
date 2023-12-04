package dataunit

import (
	"encoding/binary"
	"io"
)

// ReadDataUnit reads a single EPP data unit from [io.Reader] r, returning the payload or an error.
// If the EPP data unit read has zero-length, ReadDataUnit will return (nil, nil).
//
// An EPP data unit is prefixed with 32-bit, big-endian value specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See https://datatracker.ietf.org/doc/rfc5734/ for more information.
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
	// FIXME: do we need to support zero-length data units?
	if n == 0 {
		return nil, nil
	}
	p := make([]byte, n)
	_, err = io.ReadAtLeast(r, p, int(n))
	return p, err
}

// WriteDataUnit writes a single EPP data unit to [io.Writer] w.
//
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See https://datatracker.ietf.org/doc/rfc5734/ for more information.
func WriteDataUnit(w io.Writer, p []byte) error {
	s := uint32(4 + len(p))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	// FIXME: do we need to support zero-length data units?
	if len(p) == 0 {
		return nil
	}
	_, err = w.Write(p)
	return err
}
