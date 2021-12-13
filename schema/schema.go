package schema

// Schema represents an XML schema identified by one or more XML namespace URIs.
//
// A Schema is used to negotiate XML protocol support and allocate new instances
// of Go types associated with an xml.Name via the Factory interface.
type Schema interface {
	// SchemaName returns a short name for this Schema. For example, the EPP
	// Schema object would return "epp" and the EPP <domain> object would
	// return "domain".
	SchemaName() string

	// SchemaNS returns a slice of namespace URIs recognized by this Schema
	// in order of preference. An EPP client or server will use the first
	// matching namespace and ignore the others.
	SchemaNS() []string

	// A Schema also implements Factory.
	Factory
}
