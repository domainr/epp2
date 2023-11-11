package protocol

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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

	syncReplies chan reply

	mu           sync.Mutex
	transactions map[string]transaction
}

// NewClient returns a new EPP client using conn.
// Responses from the server will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func NewClient(conn Conn, schemas ...schema.Schema) Client {
	c := newClient(conn, schemas)
	// Read the initial <greeting> from the server.
	go c.readEPP(context.Background())
	return c
}

func newClient(conn Conn, schemas schema.Schemas) *client {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &client{
		conn:        conn,
		schemas:     schemas,
		hasGreeting: make(chan struct{}),
		syncReplies: make(chan reply),
	}
}

// Close closes the connection and cancels any pending commands.
func (c *client) Close() error {
	err := c.conn.Close()
	cerr := err
	if cerr == nil {
		cerr = ErrClosedConnection
	}
	c.cleanup(cerr)
	return err
}

// TODO: implement Shutdown(ctx) for graceful shutdown of a client connection?

// ServerConfig returns the server configuration described in a <greeting> message.
// Will block until the an initial <greeting> is received, or ctx is canceled.
//
// TODO: move this to epp.Client.
// func (c *client) ServerConfig(ctx context.Context) (Config, error) {
// 	g, err := c.Greeting(ctx)
// 	if err != nil {
// 		return Config{}, err
// 	}
// 	return configFromGreeting(g), nil
// }

// ServerName returns the most recently received server name.
// Will block until an initial <greeting> is received, or ctx is canceled.
//
// TODO: move this to epp.Client.
func (c *client) ServerName(ctx context.Context) (string, error) {
	g, err := c.Greeting(ctx)
	if err != nil {
		return "", err
	}
	return g.ServerName, nil
}

// ServerTime returns the most recently received timestamp from the server.
// Will block until an initial <greeting> is received, or ctx is canceled.
//
// TODO: move this to epp.Client.
// TODO: what is used for?
func (c *client) ServerTime(ctx context.Context) (time.Time, error) {
	g, err := c.Greeting(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return g.ServerDate.Time, nil
}

// Exchange sends [epp.Body] req and returns the response from the server.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed.
func (c *client) Exchange(ctx context.Context, req epp.Body) (epp.Body, error) {
	replies := c.syncReplies
	if cmd, ok := req.(*epp.Command); ok {
		if cmd.ClientTransactionID != "" {
			tx, cancel := newTransaction(ctx)
			defer cancel()
			err := c.pushCommand(cmd.ClientTransactionID, tx)
			if err != nil {
				return nil, err
			}
			replies = tx.reply
		}
	}

	err := c.conn.WriteEPP(req)
	if err != nil {
		return nil, err
	}

	err = c.readEPP(ctx)
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case r := <-replies:
		return r.body, r.err
	}
}

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

// readEPP reads a single EPP data unit and dispatches it to an
// awaiting transaction.
func (c *client) readEPP(ctx context.Context) error {
	body, err := c.conn.ReadEPP()
	if err != nil {
		return err
	}
	return c.handleReply(ctx, body)
}

func (c *client) handleReply(ctx context.Context, body epp.Body) error {
	replies := c.syncReplies
	switch body := body.(type) {
	case *epp.Response:
		id := body.TransactionID.Client
		if id != "" {
			tx, ok := c.popCommand(id)
			if !ok {
				// TODO: log when server responds with unknown transaction ID.
				// TODO: keep abandoned transactions around for some period of time.
				return TransactionIDError{id}
			}
			ctx = tx.ctx // TODO: should we track both contexts?
			replies = tx.reply
		}

	case *epp.Greeting:
		// Always store the last <greeting> received from the server.
		c.greeting.Store(body)

		// Close hasGreeting if this is the first <greeting> received.
		select {
		case <-c.hasGreeting:
		default:
			close(c.hasGreeting)
			// Return immediately rather than send to blocked replies channel.
			return nil
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case replies <- reply{body: body, err: nil}:
	}

	return nil
}

// pushCommand adds a <command> transaction to the map of in-flight commands.
func (c *client) pushCommand(id string, tx transaction) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.transactions[id]
	if ok {
		return DuplicateTransactionIDError{id}
	}
	if c.transactions == nil {
		c.transactions = make(map[string]transaction)
	}
	c.transactions[id] = tx
	return nil
}

// popCommand removes a <command> transaction from the map of in-flight commands.
func (c *client) popCommand(id string) (transaction, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	tx, ok := c.transactions[id]
	if ok {
		delete(c.transactions, id)
	}
	return tx, ok
}

// cleanup cleans up and responds to all in-flight <hello> and <command> transactions.
// Each transaction will be finalized with err, which may be nil.
//
// TODO: clean up stale or abandoned transactions on a regular basis?
func (c *client) cleanup(err error) {
	c.mu.Lock()
	commands := c.transactions
	c.transactions = nil
	c.mu.Unlock()
	for _, tx := range commands {
		select {
		case <-tx.ctx.Done():
		case tx.reply <- reply{err: err}:
		}
	}
}
