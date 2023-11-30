package protocol

import (
	"context"
	"sync/atomic"

	"github.com/domainr/epp2/errors"
	"github.com/domainr/epp2/schema"
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
	conn    Conn
	schemas schema.Schemas

	// greeting stores the most recently received <greeting> from the server.
	greeting atomic.Value

	// hasGreeting is closed when the client receives an initial <greeting> from the server.
	hasGreeting chan struct{}

	transactions chan transaction
}

// NewClient returns a new EPP client using conn.
// Responses from the server will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func NewClient(conn Conn, schemas ...schema.Schema) Client {
	c := newClient(conn, schemas)
	// Read the initial <greeting> from the server.
	go c.readEPP()
	return c
}

func newClient(conn Conn, schemas schema.Schemas) *client {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &client{
		conn:         conn,
		schemas:      schemas,
		hasGreeting:  make(chan struct{}),
		transactions: make(chan transaction, maxConcurrentTransactions),
	}
}

// FIXME: make this an option
const maxConcurrentTransactions = 16

// Close closes the connection and cancels any pending commands.
func (c *client) Close() error {
	err := c.conn.Close()
	cerr := err
	if cerr == nil {
		cerr = errors.ClosedConnection
	}
	c.cleanup(cerr)
	return err
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
func (c *client) Exchange(ctx context.Context, req epp.Body) (epp.Body, error) {
	// Client MUST wait until it has received a <greeting> element from the server
	// before sending any commands.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.hasGreeting:
	}

	// Queue a new transaction.
	tx, cancel := newTransaction(ctx)
	defer cancel()
	c.transactions <- tx

	// Write one request to the EPP connection.
	err := c.conn.WriteEPP(req)
	if err != nil {
		return nil, err
	}

	// Read one response off the EPP connection, which might not be a response to req.
	res, err := c.readEPP()
	if err != nil {
		return nil, err
	}

	// Dequeue a transaction.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case rtx := <-c.transactions:
		if rtx == tx {
			return res, nil // Short-circuit
		}
		// rtx.result is 1-buffered so this should never block.
		// If it blocks, then there is a bug in the implementation.
		rtx.result <- result{res, nil}
	}

	// Wait for the response to tx.
	select {
	case <-tx.ctx.Done():
		return nil, tx.ctx.Err()
	case res := <-tx.result:
		return res.body, res.err
	}
}

func (c *client) readEPP() (epp.Body, error) {
	// Read one reply, which may not be a reply to req.
	body, err := c.conn.ReadEPP()
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
	for tx := range c.transactions {
		select {
		case <-tx.ctx.Done():
		case tx.result <- result{err: err}:
		}
	}
}
