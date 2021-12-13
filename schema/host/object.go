package host

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema/epp"
)

// Object implements the epp.Object interface for the EPP <host> object type.
const Object eppObject = "host"

var _ epp.Object = Object

type eppObject string

func (o eppObject) ObjectName() string {
	return string(o)
}

func (eppObject) ObjectNS() []string {
	return []string{NS}
}

func (eppObject) New(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	// TODO: other types.
	// case "check":
	// 	return &Check{}
	}
	return nil
}
