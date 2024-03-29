package contact

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// NS defines the IETF URN for the EPP contact namespace.
// See https://www.iana.org/assignments/xml-registry/ns/contact-1.0.txt.
const NS = "urn:ietf:params:xml:ns:contact-1.0"

// Schema implements the schema.Schema interface for the EPP contact namespace.
const Schema schemaString = "contact"

var _ schema.Schema = Schema

type schemaString string

func (o schemaString) SchemaName() string {
	return string(o)
}

func (schemaString) SchemaNS() []string {
	return []string{NS}
}

func (schemaString) ResolveXML(name xml.Name) any {
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
