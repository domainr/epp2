package schema

import (
	"bytes"
	"io"

	"github.com/domainr/epp2/internal/xml"
)

// WithResolver associates a [Resolver] f with [xml.Decoder] d by overriding the
// CharsetReader field with a special value that returns the Resolver. If
// CharsetReader is not nil, the function will be wrapped.
//
// The Resolver can be retrieved from the [xml.Decoder] using GetResolver(d). This
// allows decoding of deeply-nested XML structures that are extended with types
// unknown to a parent package.
//
// If f is nil, d will not be modified.
func WithResolver(d *xml.Decoder, resolver Resolver) *xml.Decoder {
	if resolver == nil {
		return d
	}
	saved := d.CharsetReader
	d.CharsetReader = func(charset string, r io.Reader) (io.Reader, error) {
		var err error
		if saved != nil {
			r, err = saved(charset, r)
		}
		return &reader{resolver, r}, err
	}
	return d
}

// GetResolver accesses a [Resolver] associated with [xml.Decoder] d. If d does not
// have an associated Resolver, it will return nil.
func GetResolver(d *xml.Decoder) Resolver {
	if d.CharsetReader == nil {
		return nil
	}
	r, err := d.CharsetReader("utf-8", eof{})
	if err != nil {
		return nil
	}
	if resolver, ok := r.(Resolver); ok {
		return resolver
	}
	return nil
}

type eof struct{}

func (eof) Read([]byte) (int, error) {
	return 0, io.EOF
}

// UseResolver associates resolver with [xml.Decoder] d and calls f with the
// modified [xml.Decoder]. The xml.Decoder is restored before returning.
//
// If r is nil, f will be called with an unmodified xml.Decoder.
func UseResolver(d *xml.Decoder, resolver Resolver, f func(*xml.Decoder) error) error {
	saved := d.CharsetReader
	d = WithResolver(d, resolver)
	err := f(d)
	d.CharsetReader = saved
	return err
}

type reader struct {
	Resolver
	io.Reader
}

var _ io.Reader = &reader{}
var _ Resolver = &reader{}

// ResolveXML implements the [Resolver] interface.
func (r *reader) ResolveXML(name xml.Name) any {
	v := r.Resolver.ResolveXML(name)
	if v != nil {
		return v
	}

	// If r.Reader also implements [Resolver] (which means it’s probably a
	// resolverReader), call it.
	if resolver, ok := r.Reader.(Resolver); ok {
		return resolver.ResolveXML(name)
	}
	return nil
}

// Unmarshal attempts to decode p into v using [Resolver] f.
func Unmarshal(p []byte, v any, resolver Resolver) error {
	return WithResolver(xml.NewDecoder(bytes.NewReader(p)), resolver).Decode(v)
}

// DecodeElement attempts to decode start using a [Resolver] associated with d.
// Unrecognized tag names will be decoded into an instance of [Any].
func DecodeElement(d *xml.Decoder, start xml.StartElement) (any, error) {
	var v any
	r := GetResolver(d)
	if r != nil {
		v = r.ResolveXML(start.Name)
	}
	if v == nil {
		v = &Any{}
	}
	err := d.DecodeElement(v, &start)
	return v, err
}

// DecodeElements attempts to decode a sequence of XML elements using a [Resolver]
// associated with d. Unrecognized tag names will be decoded into an instance of
// [Any]. Func f will be called for each decoded element. Decoding will
// stop if f returns an error.
func DecodeElements(d *xml.Decoder, f func(any) error) error {
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if start, ok := tok.(xml.StartElement); ok {
			v, err := DecodeElement(d, start)
			if err != nil {
				return err
			}
			err = f(v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
