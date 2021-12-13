package contact

import (
	"github.com/domainr/epp2/schema/epp"
	"github.com/nbio/xml"
)

// Object implements the epp.Object interface for the EPP <contact> object type.
const Object eppObject = "contact"

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
