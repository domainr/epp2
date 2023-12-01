package schema

import "github.com/domainr/epp2/internal/xml"

// Schema represents an XML schema identified by one or more XML namespace URIs.
//
// A Schema is used to negotiate XML protocol support and allocate new instances
// of Go types associated with an xml.Name via the Resolver interface.
type Schema interface {
	// SchemaName returns a short name for this Schema. For example, the EPP
	// Schema object would return "epp" and the EPP <domain> object would
	// return "domain".
	SchemaName() string

	// SchemaNS returns a slice of namespace URIs recognized by this Schema
	// in order of preference. An EPP client or server will use the first
	// matching namespace and ignore the others.
	SchemaNS() []string

	// A Schema also implements Resolver.
	Resolver
}

// Schemas is a slice of one or more [Schema] values. It implements the
// Resolver interface, trying each Schema in order until one returns a non-nil
// value.
type Schemas []Schema

// ResolveXML tries each [Resolver] in order, returning the first non-nil value.
func (schemas Schemas) ResolveXML(name xml.Name) any {
	for _, s := range schemas {
		v := s.ResolveXML(name)
		if v != nil {
			return v
		}
	}
	return nil
}
