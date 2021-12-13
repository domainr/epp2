package domain

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Schema implements the schema.Schema interface for the EPP <domain> object type.
const Schema schemaString = "domain"

var _ schema.Schema = Schema

type schemaString string

func (s schemaString) SchemaName() string {
	return string(s)
}

func (schemaString) SchemaNS() []string {
	return []string{NS}
}

func (schemaString) New(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	// TODO: other types.
	case "check":
		return &Check{}
	}
	return nil
}
