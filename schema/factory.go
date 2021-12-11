package schema

import (
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
