// Package dataunit implements low-level encoding of the EPP data unit protocol.
//
// The data unit protocol is a simple framing of an XML payload prefixed with a 4-byte header
// in network byte order that expresses the total length of the data unit (header + payload size).
// EPP data units are sent and received via an underlying transport (typically a TLS connection).
//
// See [RFC 5730] (Extensible Provisioning Protocol) and
// [RFC 5734] (EPP Transport over TCP) for more information.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
// [RFC 5734]: https://datatracker.ietf.org/doc/rfc5734/
package dataunit

import (
	"context"
	"encoding/binary"
	"io"
)

// Receive receives a single EPP data unit from [io.Reader] r, returning the payload or an error.
// If the EPP data unit read has zero-length, Receive will return (nil, nil).
//
// It will block until a data unit is read, ctx is cancelled, or r is closed.
// The supplied Context must be non-nil.
// Concurrent calls to Receive must be protected by a synchronization mechanism.
//
// See [Read] for additional information.
func Receive(ctx context.Context, r io.Reader) ([]byte, error) {
	err := context.Cause(ctx)
	if err != nil {
		return nil, err
	}
	ch := make(chan result)
	go func() {
		data, err := Read(r)
		ch <- result{data, err}
	}()
	select {
	case <-ctx.Done():
		return nil, context.Cause(ctx)
	case res := <-ch:
		return res.data, res.err
	}
}

// Read reads a single EPP data unit from [io.Reader] r, returning the payload or an error.
// If the EPP data unit read has zero-length, Read will return (nil, nil).
//
// An EPP data unit is prefixed with 32-bit, big-endian value specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
func Read(r io.Reader) ([]byte, error) {
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
	data := make([]byte, n)
	_, err = io.ReadAtLeast(r, data, int(n))
	return data, err
}

// Send sends a single EPP data unit to [io.Writer] w.
// It will block until a data unit is written, ctx is cancelled, or w is closed.
// The supplied Context must be non-nil.
// Concurrent calls to Send must be protected by a synchronization mechanism.
//
// See [Write] for additional information.
func Send(ctx context.Context, w io.Writer, data []byte) error {
	err := context.Cause(ctx)
	if err != nil {
		return err
	}
	ch := make(chan error)
	go func() {
		ch <- Write(w, data)
	}()
	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case err := <-ch:
		return err
	}
}

// Write writes a single EPP data unit to [io.Writer] w.
//
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
func Write(w io.Writer, data []byte) error {
	s := uint32(4 + len(data))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	// FIXME: do we need to support zero-length data units?
	if len(data) == 0 {
		return nil
	}
	_, err = w.Write(data)
	return err
}
