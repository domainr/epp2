package contact

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema/epp"
)

// Object implements the epp.Object interface for the EPP <contact> object type.
// It also implements the schema.Schema and schema.Factory interfaces.
const Object eppObject = "contact"

var _ epp.Object = Object

type eppObject string

func (eppObject) EPPObject() {}

func (o eppObject) SchemaName() string {
	return string(o)
}

func (eppObject) SchemaNS() []string {
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
