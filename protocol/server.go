package protocol

import (
	"context"

	"github.com/domainr/epp2/schema/epp"
)

// Handler is the interface implemented by any type acting
// as an EPP server as defined in [RFC 5730].
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Handler interface {
	// Hello is called in response to a <hello> message, and returns the appropriate <greeting>.
	EPPHello(ctx context.Context) (*epp.Greeting, error)

	// Command accepts an EPP command and returns an EPP response.
	//
	// To correlate it with a response, cmd must have a valid, unique
	// transaction ID.
	// TODO: should it assign a transaction ID if empty?
	EPPCommand(ctx context.Context, cmd *epp.Command) (*epp.Response, error)

	// HandleBody handles an unrecognized EPP body element.
	EPPOther(ctx context.Context, body epp.Body) (epp.Body, error)

	// ServeEPP accepts an EPP request and returns an EPP response.
	// The body argument may be nil if the EPP request contained unrecognized XML.
	// In this case, ServeEPP should return an error response.
	ServeEPP(ctx context.Context, body epp.Body) (epp.Body, error)
}

// Serve services [Conn] conn with [Server] s, returning any error that
// occurs. It will block, and may be run in a separate goroutine.
//
// Serve will immediately call s.Hello() to retrieve an [epp.Greeting] to send to the client.
//
// It does not close conn.
func Serve(conn Conn, s Handler) error {
	// TODO: should Conn hold a context?
	ctx := context.Background()
	var req, res epp.Body
	var err error
	res, err = s.EPPHello(ctx)
	if err != nil {
		return err
	}
	for {
		err = conn.WriteEPP(res)
		if err != nil {
			return err
		}
		req, err = conn.ReadEPP()
		if err != nil {
			return err
		}
		switch req := req.(type) {
		case *epp.Command:
			res, err = s.EPPCommand(ctx, req)
		case *epp.Hello:
			res, err = s.EPPHello(ctx)
		default:
			// TODO: send error to client
			// TODO: epp.ErrorResponse(err)?
			res, err = s.ServeEPP(ctx, req)
		}
		if err != nil {
			return err
		}
		err = conn.WriteEPP(res)
		if err != nil {
			return err
		}
	}
}
