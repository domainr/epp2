package epp

import "github.com/domainr/epp2/schema"

// Extension represents an EPP extension. It handles negotiating the EPP
// extensions used in a client-server connection and handles translation of XML
// to Go types.
type Extension interface {
	// ExtensionNS returns a slice of namespace URIs for this Extension in
	// order of preference. An EPP client or server will use the first
	// matching namespace and ignore the others.
	ExtensionNS() []string

	// An Extension is also a schema.Factory. Each Extension is responsible
	// for mapping an xml.Name to a concrete instance of a type when
	// decoding XML.
	schema.Factory
}
