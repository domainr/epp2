package common

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// NS defines the IETF URN for the EPP common namespace.
// See https://www.iana.org/assignments/xml-registry/ns/eppcom-1.0.txt.
const NS = "urn:ietf:params:xml:ns:eppcom-1.0"

// Schema implements the schema.Schema interface for the EPP common namespace.
const Schema schemaString = "eppcom"

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
	// TODO: what are EPP common types?
	}
	return nil
}
