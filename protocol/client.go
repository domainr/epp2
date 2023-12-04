package protocol

import (
	"context"
	"io"

	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Client represents a low-level EPP client as defined in [RFC 5730].
// A Client is safe to use from multiple goroutines.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Client interface {
	// ExchangeEPP sends an EPP message and returns an EPP response.
	// It blocks until a response is received, the Context is canceled, or
	// the underlying connection is closed.
	ExchangeEPP(context.Context, epp.Body) (epp.Body, error)
}

type client struct {
	client dataunit.Client
	coder  coder
}

// Connect connects to an EPP server over conn. It waits for the initial
// <greeting> message from the server before returning, or until ctx is cancelled
// or the underlying connection is closed.
// Responses from the server will be decoded using [schema.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func Connect(ctx context.Context, conn io.ReadWriter, schemas ...schema.Schema) (Client, epp.Body, error) {
	c := newClient(conn, schemas)

	// Read the initial <greeting> from the server.
	data, err := dataunit.Receive(ctx, conn)
	if err != nil {
		return c, nil, err
	}
	body, err := c.coder.umarshalXML(data)
	return c, body, err
}

func newClient(conn io.ReadWriter, schemas schema.Schemas) *client {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &client{
		client: dataunit.Client{Conn: conn},
		coder:  coder{schemas},
	}
}

// ExchangeEPP sends [epp.Body] req and returns the response from the server.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed.
// Exchange is safe to call from multiple goroutines.
func (c *client) ExchangeEPP(ctx context.Context, req epp.Body) (epp.Body, error) {
	data, err := c.coder.marshalXML(req)
	if err != nil {
		return nil, err
	}
	data, err = c.client.ExchangeDataUnit(ctx, data)
	if err != nil {
		return nil, err
	}
	return c.coder.umarshalXML(data)
}
