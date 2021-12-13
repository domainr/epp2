package epp

import "github.com/domainr/epp2/schema"

// Object is a generic EPP object type.
// Standard EPP objects are <domain>, <host>, and <contact>.
type Object interface {
	EPPObject()
	schema.Schema
}
