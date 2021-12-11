package schema

import (
	"io"

	"github.com/nbio/xml"
)

// Factory is a generic interface that can return a new instance of a type
// identified by an xml.Name.
//
// New should return nil for any xml.Name it does not recognize.
type Factory interface {
	New(name xml.Name) interface{}
}

// FactoryFunc is a function that implements the Factory interface.
type FactoryFunc func(name xml.Name) interface{}

// New calls f and returns the value.
func (f FactoryFunc) New(name xml.Name) interface{} {
	return f(name)
}

// Factories is a slice of one or more Factory instances. It implements the
// Factory interface, trying each Factory in order until one returns a non-nil
// value.
type Factories []Factory

var _ Factory = Factories{}

// New tries each Factory in order, returning the first non-nil value.
func (factories Factories) New(name xml.Name) interface{} {
	for _, f := range factories {
		v := f.New(name)
		if v != nil {
			return v
		}
	}
	return nil
}

// WithFactory associates a Factory f with xml.Decoder d. The Factory can be
// retrieved from the Decoder using GetFactory(d).
//
// WithFactory allows decoding of deeply-nested XML structures that are extended
// with types unknown to a parent package.
func WithFactory(d *xml.Decoder, f Factory, cb func(*xml.Decoder) error) error {
	saved := d.CharsetReader
	d.CharsetReader = (&crowbar{d.CharsetReader, f}).CharsetReader
	err := cb(d)
	d.CharsetReader = saved
	return err
}

// GetFactory accesses a Factory associated with xml.Decoder d. If d does not
// have an associated Factory, it will return nil.
func GetFactory(d *xml.Decoder) Factory {
	if d.CharsetReader == nil {
		return nil
	}
	err, r := d.CharsetReader("", nil)
	if err != nil {
		return nil
	}
	if f, ok := r.(Factory); ok {
		return f
	}
	return nil
}

type crowbar struct {
	cr func(charset string, input io.Reader) (io.Reader, error)
	f  Factory
}

// CharsetReader, if called with a nil io.Reader, will bypass the underlying
// CharsetReader func and return a factoryReader.
func (c *crowbar) CharsetReader(charset string, r io.Reader) (io.Reader, error) {
	var err error
	if c.cr != nil && r != nil {
		r, err = c.cr(charset, r)
	}
	return &factoryReader{
		Reader:  r,
		Factory: c.f,
	}, err
}

type factoryReader struct {
	Factory
	io.Reader
}

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
