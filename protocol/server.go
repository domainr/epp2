package protocol

import (
	"context"
	"io"

	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Responder is the interface implemented by any type that can respond to a client request with an EPP message body.
type Responder interface {
	// RespondEPP writes an EPP response to the client. It blocks until the message is written,
	// Context is canceled, or the underlying connection is closed.
	//
	// RespondEPP can be called from an arbriary goroutine, but must only be called once.
	RespondEPP(context.Context, epp.Body) error
}

type responderFunc func(context.Context, epp.Body) error

func (f responderFunc) RespondEPP(ctx context.Context, body epp.Body) error {
	return f(ctx, body)
}

// Server represents a a low-level EPP server connection, as defined in [RFC 5730].
// A Server is safe to use from multiple goroutines.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Server interface {
	// ServeEPP provides an EPP request and a mechanism to respond to the request.
	// It blocks until a response is received, Context is canceled, or the underlying connection is closed.
	//
	// The supplied Context must be non-nil, and only affects reading the request from the client.
	// Cancelling the Context after ServeEPP returns will have no effect on the Responder.
	//
	// The returned [Responder] should only be used once. The returned Responder will always
	// be non-nil, so the caller can respond to a malformed client request.
	ServeEPP(context.Context) (epp.Body, Responder, error)
}

type server struct {
	server dataunit.Server
	coder  coder
}

// Serve services conn as an EPP server, sending the initial <greeting> to the client.
//
// The supplied Context will be used only for sending the initial greeting.
// Cancelling ctx after Serve returns will have no effect on the resulting connection.
//
// EPP requests from the client will be decoded using [schema.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func Serve(ctx context.Context, conn io.ReadWriter, greeting epp.Body, schemas ...schema.Schema) (Server, error) {
	s := newServer(conn, schemas)
	// Send the initial <greeting> to the client.
	data, err := s.coder.marshal(greeting)
	if err != nil {
		return nil, err
	}
	return s, dataunit.Send(ctx, conn, data)
}

func newServer(conn io.ReadWriter, schemas schema.Schemas) *server {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &server{
		server: dataunit.Server{Conn: conn},
		coder:  coder{schemas},
	}
}

func (s *server) ServeEPP(ctx context.Context) (epp.Body, Responder, error) {
	data, r, err := s.server.ServeDataUnit(ctx)
	f := responderFunc(func(ctx context.Context, body epp.Body) error {
		data, err := s.coder.marshal(body)
		if err != nil {
			return err
		}
		return r.RespondDataUnit(ctx, data)
	})
	if err != nil {
		return nil, f, err
	}
	body, err := s.coder.unmarshal(data)
	return body, f, err
}
