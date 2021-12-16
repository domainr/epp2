package epp

import (
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// builder is a documentation interface to describe an ideal Builder type.
type builder interface {
	ForLogin(*epp.Login) (*Builder, error)
	ForGreeting(*epp.Greeting) (*Builder, error)

	Greeting() (*epp.Greeting, error)
	Command(action epp.Action, extensions ...epp.Extension) (*epp.Command, error)
	Login(clientID, password string, newPassword *string) (*epp.Command, error)
	Logout() (*epp.Command, error)
	ErrorResponse(error) *epp.Response
}

// Builder builds EPP elements suitable for marshaling into XML. It abstracts
// away EPP object and extension configuration and can negotiate supported and
// announced EPP namespaces between a client and server. The objects created by
// a Builder can be transmitted via the protocol package.
type Builder struct {
	// Config is the protocol configuration for this Builder. This includes
	// which protocol version to use (default: 1.0), language(s) to use
	// (default: "en"), EPP objects, and EPP extensions.
	//
	// If nil, or if individual fields in Config are nil or empty,
	// reasonable defaults will be used. If Config.Objects or
	// Config.Extensions are nil or empty, the namespaces from Schemas will
	// be used.
	Config *Config

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

// ForGreeting returns a client-centric Builder sharing the mutual capabilities announced in an
// epp.Greeting. If the Builder or Greeting do not share a mutual set of EPP versions,
// languages, and objects, it will return an error. The resulting Builder is
// suitable for creating EPP elements transmittable to the server that sent the Greeting.
func (b *Builder) ForGreeting(greeting *epp.Greeting) (*Builder, error) {
	return b, nil // TODO
}

// ForLogin returns a server-centric Builder sharing the mutual capabilities announced in an epp.Login.
// If the Builder or Login do not share a mutual set of EPP versions, languages,
// and objects, it will return an error. The resulting Builder is
// suitable for creating EPP elements transmittable to the client that generated the Login.
//
// The error returned may be an *epp.Result which can be transmitted back to an
// EPP client in an epp.Response.
func (b *Builder) ForLogin(login *epp.Login) (*Builder, error) {
	return b, nil // TODO
}

func (b *Builder) Greeting() (*epp.Greeting, error) {
	return &epp.Greeting{}, nil // TODO
}

func (b *Builder) Command(action epp.Action, extensions ...epp.Extension) (*epp.Command, error) {
	return &epp.Command{
		Action:              action,
		Extensions:          extensions,
		ClientTransactionID: b.TransactionID(),
	}, nil
}

func (b *Builder) Login(clientID, password string, newPassword *string) (*epp.Command, error) {
	return b.Command(&epp.Login{
		ClientID:    clientID,
		Password:    password,
		NewPassword: newPassword,
		Options: epp.Options{
			Version: epp.Version,
		},
	})
}

func (b *Builder) Logout() (*epp.Command, error) {
	return b.Command(&epp.Logout{})
}
