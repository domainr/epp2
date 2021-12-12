package schema

import (
	"github.com/nbio/xml"
)

// Factory is a generic interface that can return a new instance of a type
// identified by an xml.Name.
//
// New must return nil for any xml.Name it does not recognize.
type Factory interface {
	New(name xml.Name) interface{}
}

// FactoryFunc is a function that implements the Factory interface.
type FactoryFunc func(name xml.Name) interface{}

// New calls f and returns the value.
func (f FactoryFunc) New(name xml.Name) interface{} {
	return f(name)
}

// Factories merges multiple Factory instances into a single Factory. It
// implements the Factory interface, trying each Factory in order from first to
// last until one returns a non-nil value.
func Factories(f ...Factory) Factory {
	return flatten(factories(f))
}

func flatten(in factories) factories {
	if len(in) == 0 {
		return in
	}
	var out factories
	for _, f := range in {
		if f == nil {
			continue
		}
		if slice, ok := f.(factories); ok {
			out = append(out, flatten(slice)...)
		} else {
			out = append(out, f)
		}
	}
	return out
}

// factories is a slice of one or more Factory instances. It implements the
// Factory interface, trying each Factory in order until one returns a non-nil
// value.
type factories []Factory

// New tries each Factory in order, returning the first non-nil value.
func (slice factories) New(name xml.Name) interface{} {
	for _, f := range slice {
		v := f.New(name)
		if v != nil {
			return v
		}
	}
	return nil
}
