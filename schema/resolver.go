package schema

import (
	"github.com/domainr/epp2/internal/xml"
)

// Resolver is the interface implemented by any type can can resolve
// an [xml.Name] into some value.
//
// ResolveXML must return nil for any xml.Name it does not recognize.
type Resolver interface {
	ResolveXML(name xml.Name) any
}

// ResolverFunc is a func that implements the [Resolver] interface.
type ResolverFunc func(name xml.Name) any

// ResolveXML implements [Resolver] via func f.
func (f ResolverFunc) ResolveXML(name xml.Name) any {
	return f(name)
}

// Flatten merges multiple [Resolver] instances together into a single Resolver. It
// implements the Resolver interface, trying each Resolver in order from first to
// last until one returns a non-nil value.
func Flatten(f ...Resolver) Resolver {
	return flatten(resolvers(f))
}

func flatten(in resolvers) resolvers {
	if len(in) == 0 {
		return in
	}
	var out resolvers
	for _, f := range in {
		if f == nil {
			continue
		}
		if slice, ok := f.(resolvers); ok {
			out = append(out, flatten(slice)...)
		} else {
			out = append(out, f)
		}
	}
	return out
}

// resolvers is a slice of one or more [Resolver] instances. It implements the
// Resolver interface, trying each Resolver in order until one returns a non-nil
// value.
type resolvers []Resolver

// ResolveXML tries each [Resolver] in order, returning the first non-nil value.
func (slice resolvers) ResolveXML(name xml.Name) any {
	for _, f := range slice {
		v := f.ResolveXML(name)
		if v != nil {
			return v
		}
	}
	return nil
}
