package schema

import (
	"bytes"
	"io"

	"github.com/nbio/xml"
)

// WithFactory associates a Factory f with xml.Decoder d by overriding the
// CharsetReader field with a special value that returns the Factory. If
// CharsetReader is not nil, the function will be wrapped.
//
// The Factory can be retrieved from the Decoder using GetFactory(d). This
// allows decoding of deeply-nested XML structures that are extended with types
// unknown to a parent package.
//
// If f is nil, d will not be modified.
func WithFactory(d *xml.Decoder, f Factory) *xml.Decoder {
	if f == nil {
		return d
	}
	saved := d.CharsetReader
	d.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		var err error
		if saved != nil {
			r, err = saved(charset, r)
		}
		return &factoryReader{f, r}, err
	}
	return d
}

// GetFactory accesses a Factory associated with xml.Decoder d. If d does not
// have an associated Factory, it will return nil.
func GetFactory(d *xml.Decoder) Factory {
	if d.CharsetReader == nil {
		return nil
	}
	r, err := d.CharsetReader("utf-8", eof{})
	if err != nil {
		return nil
	}
	if f, ok := r.(Factory); ok {
		return f
	}
	return nil
}

type eof struct{}

func (eof) Read([]byte) (int, error) {
	return 0, io.EOF
}

// UseFactory associates a Factory f with xml.Decoder d and calls cb with the
// modified Decoder. It restores the Decoder before returning.
//
// If f is nil, cb will be called with an unmodified xml.Decoder.
func UseFactory(d *xml.Decoder, f Factory, cb func(*xml.Decoder) error) error {
	saved := d.CharsetReader
	d = WithFactory(d, f)
	err := cb(d)
	d.CharsetReader = saved
	return err
}

type factoryReader struct {
	Factory
	io.Reader
}

var _ io.Reader = &factoryReader{}
var _ Factory = &factoryReader{}

// New implements the Factory interface.
func (r *factoryReader) New(name xml.Name) interface{} {
	v := r.Factory.New(name)
	if v != nil {
		return v
	}

	// If r.Reader also implements Factory (which means it’s probably a
	// factoryReader), call it.
	if f, ok := r.Reader.(Factory); ok {
		return f.New(name)
	}
	return nil
}

// Unmarshal attempts to decode p into v using Factory f.
func Unmarshal(p []byte, v interface{}, f Factory) error {
	return WithFactory(xml.NewDecoder(bytes.NewReader(p)), f).Decode(v)
}

// DecodeElement attempts to decode start using a Factory associated with d.
// Unrecognized tag names will be decoded into an instance of Any.
func DecodeElement(d *xml.Decoder, start xml.StartElement) (interface{}, error) {
	var v interface{}
	f := GetFactory(d)
	if f != nil {
		v = f.New(start.Name)
	}
	if v == nil {
		v = &Any{}
	}
	err := d.DecodeElement(v, &start)
	return v, err
}

// DecodeElements will read and attempt to decode a sequence of XML elements
// using a Factory associated with d. Unrecognized tag names will be decoded
// into an instance of Any.
func DecodeElements(d *xml.Decoder) ([]interface{}, error) {
	var values []interface{}
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return values, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			v, err := DecodeElement(d, start)
			if err != nil {
				return values, err
			}
			values = append(values, v)
		}
	}
	return values, nil
}
