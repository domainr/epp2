package protocol

import (
	"encoding/binary"
	"io"
)

// ReadDataUnit reads a single EPP data unit from r, returning the payload or an error.
// An EPP data unit is prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
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

// WriteDataUnit writes a single EPP data unit to w.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func WriteDataUnit(w io.Writer, p []byte) error {
	s := uint32(4 + len(p))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(p)
	return err
}
