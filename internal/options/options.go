package options

import "github.com/domainr/epp2/internal"

type Options interface {
	EPPOptions(internal.Internal)
}

// Struct is an optimized form of EPP options,
// suitable for passing via the call stack.
type Struct struct {
}

func GetOption[T any](opts Options, setter func(T) Options) (T, bool) {

}
