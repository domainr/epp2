package host

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// Schema implements the schema.Schema interface for the EPP host object.
const Schema schemaString = "host"

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
	// case "check":
	// 	return &Check{}
	}
	return nil
}
