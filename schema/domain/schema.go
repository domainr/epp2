package domain

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// NS defines the IETF URN for the EPP domain namespace.
// See https://www.iana.org/assignments/xml-registry/ns/domain-1.0.txt
// and https://datatracker.ietf.org/doc/html/rfc5731.
const NS = "urn:ietf:params:xml:ns:domain-1.0"

// Schema implements the schema.Schema interface for the EPP domain namespace.
const Schema schemaString = "domain"

var _ schema.Schema = Schema

type schemaString string

func (s schemaString) SchemaName() string {
	return string(s)
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
	case "check":
		return &Check{}
	}
	return nil
}
