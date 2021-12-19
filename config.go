package epp

import (
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Config describes the configuration of an EPP client or server, including EPP
// objects and extensions used for a connection.
type Config struct {
	// Supported EPP version(s). Typically this should not be set by either
	// a client or server. If nil, this will default to []string{"1.0"}
	// (currently the only supported version).
	Versions []string

	// BCP 47 language code(s) for human-readable messages. For clients,
	// this describes the desired language(s) in preferred order. If the
	// server does not support any of the client’s preferred languages, the
	// first language advertised by the server will be selected. For
	// servers, this describes its supported language(s). If nil,
	// []string{"en"} will be used.
	Languages []string

	// Objects is a list of XML namespace URIs enumerating the EPP objects
	// supported by the client or server.
	//
	// For clients, this describes the object type(s) the client wants to
	// access.
	//
	// For servers, this describes the object type(s) the server allows
	// clients to access.
	//
	// If nil, a reasonable set of defaults will be used.
	Objects []string

	// Extensions is a list of XML namespace URIs enumerating the EPP
	// extensions supported by the client or server.
	//
	// For clients, this is a list of extensions(s) the client wants to use
	// in preferred order. If nil, a client will use its preferred version
	// of each supported extension announced by the server.
	//
	// For servers, this is list of announced extension(s).
	//
	// If nil, no EPP extensions will be used.
	Extensions []string

	// UnannouncedExtensions contains one or more EPP extensions to be used even if
	// the peer does not announce support for it. This is used as a
	// workaround for EPP servers that incorrectly announce the extensions
	// they support.
	//
	// This value should typically be left nil. This will always be nil when
	// read from a peer.
	UnannouncedExtensions []string

	// Schemas contains the Schema objects used to map object and extension
	// namespaces to Go types.
	//
	// If nil or empty, reasonable defaults will be used.
	Schemas []schema.Schema

	// TransactionID, if not nil, returns unique values used for client or
	// server transaction IDs. For clients, this generates command
	// transaction IDs. For servers, this generates response transaction
	// IDs.
	//
	// The function must be safe to call from multiple goroutines.
	//
	// If nil, a sequential transaction ID with a random prefix will be
	// used.
	TransactionID func() string
}

// Copy deep copy of c.
func (c Config) Copy() Config {
	c.Versions = copySlice(c.Versions)
	c.Languages = copySlice(c.Languages)
	c.Objects = copySlice(c.Objects)
	c.Extensions = copySlice(c.Extensions)
	c.UnannouncedExtensions = copySlice(c.UnannouncedExtensions)
	return c
}

func copySlice(s []string) []string {
	if s == nil {
		return nil
	}
	dst := make([]string, len(s))
	copy(dst, s)
	return dst
}

// ConfigForGreeting returns a client-centric Config sharing the mutual
// capabilities announced in an epp.Greeting. If the Config or Greeting do not
// share a mutual set of EPP versions, languages, and objects, it will return an
// error. The resulting Config is suitable for creating EPP elements
// transmittable to the server that sent the Greeting.
//
// TODO: implement this function.
func ConfigForGreeting(cfg *Config, greeting *epp.Greeting) (*Config, error) {
	return cfg, nil // TODO
}

// ConfigForLogin returns a server-centric Config sharing the mutual
// capabilities announced in an epp.Login. If the Config or Login do not share a
// mutual set of EPP versions, languages, and objects, it will return an error.
// The resulting Config is suitable for creating EPP elements transmittable to
// the client that generated the Login.
//
// The error returned may be an *epp.Result which can be transmitted back to an
// EPP client in an epp.Response.
//
// TODO: implement this function.
func ConfigForLogin(cfg *Config, login *epp.Login) (*Config, error) {
	return cfg, nil // TODO
}
