package epp

import "github.com/domainr/epp2/schema"

// Object is a generic EPP object type. The standard EPP object types are
// <domain>, <host>, and <contact>.
type Object interface {
	// ObjectName returns a short name for this Object.
	// For example, the EPP <domain> object would return "domain".
	ObjectName() string

	// ObjectNS returns a slice of namespace URIs for this Object in order
	// of preference. An EPP client or server will use the first matching
	// namespace and ignore the others.
	ObjectNS() []string

	// An Object is also a schema.Factory. Each Object is responsible for
	// mapping an xml.Name to a concrete instance of a type when decoding
	// XML.
	schema.Factory
}
