package domain

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Factory returns a schema.Factory for types defined in the this package.
func Factory() schema.Factory {
	return schema.FactoryFunc(factory)
}

func factory(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	case "check":
		return &Check{}
	}
	return nil
}
