package epp

import (
	"github.com/domainr/epp2/internal/xml"

	"github.com/domainr/epp2/schema"
)

// NS is the IETF URN for the EPP namespace.
// See https://www.iana.org/assignments/xml-registry/ns/epp-1.0.txt.
const NS = "urn:ietf:params:xml:ns:epp-1.0"

// Schema implements the schema.Schema interface for the the core EPP namespace.
const Schema schemaString = "epp"

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
	// Body
	case "hello":
		return &Hello{}
	case "greeting":
		return &Greeting{}
	case "command":
		return &Command{}
	case "response":
		return &Response{}

	// CommandType
	case "check":
		return &Check{}
	case "create":
		return &Create{}
	case "delete":
		return &Delete{}
	case "info":
		return &Info{}
	case "login":
		return &Login{}
	case "logout":
		return &Logout{}
	case "poll":
		return &Poll{}
	case "renew":
		return &Renew{}
	case "transfer":
		return &Transfer{}
	case "update":
		return &Update{}
	}
	return nil
}
