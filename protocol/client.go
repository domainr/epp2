package protocol

import (
	"context"

	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Client is a low-level client for the Extensible Provisioning Protocol (EPP)
// as defined in [RFC 5730]. A Client is safe to use from multiple goroutines.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Client interface {
	// ExchangeEPP sends an EPP message and returns an EPP response.
	// It blocks until a response is received, the Context is canceled, or
	// the underlying connection is closed.
	ExchangeEPP(context.Context, epp.Body) (epp.Body, error)

	// Close closes the connection.
	Close() error
}

type client struct {
	client dataunit.Client
	coder  coder
}

// Connect connects to an EPP server over conn. It waits for the initial
// <greeting> message from the server before returning, or until ctx is cancelled
// or the underlying connection is closed.
// Responses from the server will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func Connect(ctx context.Context, conn dataunit.Conn, schemas ...schema.Schema) (Client, epp.Body, error) {
	c := newClient(conn, schemas)

	// Read the initial <greeting> from the server.
	ch := make(chan result, 1)
	go func() {
		data, err := conn.ReadDataUnit()
		ch <- result{data, err}
	}()

	select {
	case <-ctx.Done():
		return c, nil, context.Cause(ctx)
	case res := <-ch:
		body, err := c.coder.umarshalXML(res.data)
		return c, body, err
	}
}

type result struct {
	data []byte
	err  error
}

func newClient(conn dataunit.Conn, schemas schema.Schemas) *client {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &client{
		client: dataunit.Client{Conn: conn},
		coder:  coder{schemas},
	}
}

// Close closes the connection, interrupting any in-flight requests.
func (c *client) Close() error {
	return c.client.Conn.Close()
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
