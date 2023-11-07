package contact

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// Schema implements the schema.Schema interface for the EPP contact object.
const Schema schemaString = "contact"

var _ schema.Schema = Schema

type schemaString string

func (o schemaString) SchemaName() string {
	return string(o)
}

func (schemaString) SchemaNS() []string {
	return []string{NS}
}

func (schemaString) New(name xml.Name) any {
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
