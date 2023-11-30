package protocol

import (
	"context"
	"sync/atomic"

	"github.com/domainr/epp2/internal/config"
	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema/epp"
)

// Client is a low-level client for the Extensible Provisioning Protocol (EPP)
// as defined in [RFC 5730]. A Client is safe to use from multiple goroutines.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Client interface {
	// Exchange sends an EPP message and returns an EPP response.
	// It blocks until a response is received, ctx is canceled, or
	// the underlying connection is closed.
	Exchange(context.Context, epp.Body) (epp.Body, error)

	// Greeting returns the last <greeting> received from the server.
	// It blocks until the <greeting> is received, ctx is canceled, or
	// the underlying connection is closed.
	Greeting(context.Context) (*epp.Greeting, error)

	// Close closes the connection.
	Close() error
}

type client struct {
	client dataunit.Client
	coder  coder

	ctx    context.Context
	cancel func()

	// hasGreeting is closed when the client receives an initial <greeting> from the server.
	hasGreeting chan struct{}

	// greeting stores the most recently received <greeting> from the server.
	greeting atomic.Value

	// err stores the first non-recoverable error on the connection.
	err atomic.Value
}

// NewClient returns a new EPP client using conn.
// Responses from the server will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func NewClient(conn dataunit.Conn, opts ...Options) Client {
	c := newClient(conn, opts...)
	// Read the initial <greeting> from the server.
	go func() {
		data, err := conn.ReadDataUnit()
		if err != nil {
			c.err.Store(err)
			c.cancel()
			return
		}
		_, err = c.receiveDataUnit(data)
		if err != nil {
			c.err.Store(err)
			c.cancel()
			return
		}
	}()
	return c
}

func newClient(conn dataunit.Conn, opts ...Options) *client {
	var cfg config.Config
	cfg.Join(opts...)

	if cfg.Context == nil {
		cfg.Context = context.Background()
	}
	if len(cfg.Schemas) == 0 {
		cfg.Schemas = DefaultSchemas()
	}

	ctx, cancel := context.WithCancel(cfg.Context)

	return &client{
		client:      dataunit.NewClient(conn),
		coder:       coder{cfg.Schemas},
		ctx:         ctx,
		cancel:      cancel,
		hasGreeting: make(chan struct{}),
	}
}

// Close closes the connection and cancels any pending commands.
func (c *client) Close() error {
	// err := c.conn.Close()
	// cerr := err
	// if cerr == nil {
	// 	cerr = errors.ClosedConnection
	// }
	// c.cleanup(cerr)
	// return err
	return nil
}

// TODO: implement Shutdown(ctx) for graceful shutdown of a client connection?

// Greeting returns the last <greeting> received from the server.
// It blocks until the <greeting> is received, ctx is canceled, or
// the underlying connection is closed.
func (c *client) Greeting(ctx context.Context) (*epp.Greeting, error) {
	g := c.greeting.Load()
	if g != nil {
		return g.(*epp.Greeting), nil
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.hasGreeting:
		return c.greeting.Load().(*epp.Greeting), nil
	}
}

// Exchange sends [epp.Body] req and returns the response from the server.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed.
// Exchange is safe to call from multiple goroutines.
func (c *client) Exchange(ctx context.Context, req epp.Body) (epp.Body, error) {
	// Client MUST wait until it has received a <greeting> element from the server
	// before sending any commands.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.hasGreeting:
	}

	reqData, err := c.coder.marshalXML(req)
	if err != nil {
		return nil, err
	}

	ch := make(chan result, 1)

	go func() {
		resData, err := c.client.SendDataUnit(reqData)
		ch <- result{resData, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case res := <-ch:
		if res.err != nil {
			return nil, err
		}
		return c.receiveDataUnit(res.data)
	}
}

type result struct {
	data []byte
	err  error
}

func (c *client) receiveDataUnit(data []byte) (epp.Body, error) {
	body, err := c.coder.umarshalXML(data)
	if err != nil {
		return nil, err
	}

	switch body := body.(type) {
	case *epp.Greeting:
		// Always store the last <greeting> received from the server.
		c.greeting.Store(body)

		// Close hasGreeting if this is the first <greeting> received.
		select {
		case <-c.hasGreeting:
		default:
			close(c.hasGreeting)
		}
	}

	return body, err
}

// cleanup cleans up and responds to all in-flight transactions.
// Each transaction will be finalized with err, which may be nil.
//
// TODO: clean up stale or abandoned transactions on a regular basis?
func (c *client) cleanup(err error) {
	// TODO: set an error on c with an atomic.Value so in-flight transactions can exit gracefully.
}
