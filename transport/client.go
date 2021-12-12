package transport

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema/epp"
)

// Client is a low-level client for the Extensible Provisioning Protocol (EPP)
// as defined in RFC 3790. See https://www.rfc-editor.org/rfc/rfc5730.html.
// A Client is safe to use from multiple goroutines.
type Client interface {
	// Command sends an EPP command and returns an EPP response.
	// It blocks until a response is received, ctx is canceled, or
	// the underlying connection is closed.
	//
	// To correlate it with a response, cmd must have a valid, unique
	// transaction ID.
	// TODO: should it assign a transaction ID if empty?
	Command(ctx context.Context, cmd *epp.Command) (*epp.Response, error)

	// Hello sends an EPP <hello> and returns the <greeting> received.
	// It blocks until a <greeting> is received, ctx is canceled, or
	// the underlying connection is closed.
	Hello(ctx context.Context) (*epp.Greeting, error)

	// Greeting returns the last <greeting> recieved from the server.
	// It blocks until the <greeting> is received, ctx is canceled, or
	// the underlying connection is closed.
	Greeting(ctx context.Context) (*epp.Greeting, error)

	// Close closes the connection.
	Close() error
}

type client struct {
	// reading synchronizes reads from transport.
	reading sync.Mutex

	// writing synchronizes writes to transport.
	writing sync.Mutex

	// TODO: rename this.
	transport Transport

	// greeting stores the most recently received <greeting> from the server.
	greeting atomic.Value

	// hasGreeting is closed when the client receives an initial <greeting> from the server.
	hasGreeting chan struct{}

	mHellos sync.Mutex
	hellos  []transaction

	mCommands sync.Mutex
	commands  map[string]transaction

	// done is closed when the client receives a fatal error or the connection is closed.
	done chan struct{}
}

// NewClient returns a new Transport using conn.
func NewClient(conn Transport) Client {
	c := newClient(conn)
	// Read the initial <greeting> from the server.
	go c.readEPP()
	return c
}

func newClient(t Transport) *client {
	return &client{
		transport:   t,
		hasGreeting: make(chan struct{}),
		commands:    make(map[string]transaction),
		done:        make(chan struct{}),
	}
}

// Close closes the connection and cancels any pending commands.
func (c *client) Close() error {
	err := c.transport.Close()
	cerr := err
	if cerr == nil {
		cerr = ErrClosedConnection
	}
	c.cleanup(cerr)
	return err
}

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

// Command sends an EPP command and returns an EPP response.
// It blocks until a response is received, ctx is canceled, or
// the underlying connection is closed.
func (c *client) Command(ctx context.Context, cmd *epp.Command) (*epp.Response, error) {
	tx, cancel := newTransaction(ctx)
	defer cancel()
	c.pushCommand(cmd.ClientTransactionID, tx)

	err := c.writeEPP(cmd)
	if err != nil {
		return nil, err
	}

	err = c.readEPP()
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case reply := <-tx.reply:
		if r, ok := reply.body.(*epp.Response); ok {
			return r, reply.err
		}
		return nil, reply.err
	}
}

// Hello sends an EPP <hello> message to the server.
// It will block until the next <greeting> message is received or ctx is canceled.
func (c *client) Hello(ctx context.Context) (*epp.Greeting, error) {
	tx, cancel := newTransaction(ctx)
	defer cancel()
	c.pushHello(tx)

	err := c.writeEPP(&epp.Hello{})
	if err != nil {
		return nil, err
	}

	err = c.readEPP()
	if err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case reply := <-tx.reply:
		if g, ok := reply.body.(*epp.Greeting); ok {
			return g, reply.err
		}
		return nil, reply.err
	}
}

// Greeting returns the last <greeting> recieved from the server.
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

// writeEPP writes body to the underlying Transport.
// Writes are synchronized, so it is safe to call this from multiple goroutines.
func (c *client) writeEPP(body epp.Body) error {
	x, err := xml.Marshal(epp.EPP{Body: body})
	if err != nil {
		return err
	}
	return c.writeDataUnit(x)
}

// writeDataUnit writes a single EPP data unit to the underlying Transport.
// Writes are synchronized, so it is safe to call this from multiple goroutines.
func (c *client) writeDataUnit(p []byte) error {
	c.writing.Lock()
	err := c.transport.WriteDataUnit(p)
	c.writing.Unlock()
	return err
}

// readEPP reads a single EPP data unit from c.t and dispatches it to an
// awaiting transaction.
func (c *client) readEPP() error {
	c.reading.Lock()
	p, err := c.transport.ReadDataUnit()
	c.reading.Unlock()
	if err != nil {
		return err
	}
	return c.handleDataUnit(p)
}

func (c *client) handleDataUnit(p []byte) error {
	var e epp.EPP
	err := xml.Unmarshal(p, &e)
	if err != nil {
		return err
	}
	return c.handleReply(e.Body)
}

func (c *client) handleReply(body epp.Body) error {
	switch body := body.(type) {
	case *epp.Response:
		id := body.TransactionID.Client
		if id == "" {
			// TODO: log when server responds with an empty client transaction ID.
			return TransactionIDError(id)
		}
		t, ok := c.popCommand(id)
		if !ok {
			// TODO: log when server responds with unknown transaction ID.
			// TODO: keep abandoned transactions around for some period of time.
			return TransactionIDError(id)
		}
		err := c.replyTo(t, body, nil)
		if err != nil {
			return err
		}

	case *epp.Greeting:
		// Always store the last <greeting> received from the server.
		c.greeting.Store(body)

		// Close hasGreeting this is the first <greeting> recieved.
		select {
		case <-c.hasGreeting:
		default:
			close(c.hasGreeting)
		}

		// Pass the <greeting> to a caller waiting on it.
		t, ok := c.popHello()
		if ok {
			err := c.replyTo(t, body, nil)
			if err != nil {
				return err
			}
		}

	case *epp.Hello:
		// TODO: log if server receives a <hello> or <command>.

	case *epp.Command:
		// TODO: log if server receives a <hello> or <command>.
	}

	return nil
}

func (c *client) replyTo(t transaction, body epp.Body, err error) error {
	select {
	case <-t.ctx.Done():
		return t.ctx.Err()
	case t.reply <- reply{body: body, err: err}:
	}
	return nil
}

// pushHello adds a <hello> transaction to the end of the stack.
func (c *client) pushHello(tx transaction) {
	c.mHellos.Lock()
	defer c.mHellos.Unlock()
	c.hellos = append(c.hellos, tx)
}

// popHello pops the oldest <hello> transaction off the front of the stack.
func (c *client) popHello() (transaction, bool) {
	c.mHellos.Lock()
	defer c.mHellos.Unlock()
	if len(c.hellos) == 0 {
		return transaction{}, false
	}
	tx := c.hellos[0]
	c.hellos = c.hellos[1:]
	return tx, true
}

// pushCommand adds a <command> transaction to the map of in-flight commands.
func (c *client) pushCommand(id string, tx transaction) error {
	c.mCommands.Lock()
	defer c.mCommands.Unlock()
	_, ok := c.commands[id]
	if ok {
		return fmt.Errorf("epp: transaction already exists: %s", id)
	}
	c.commands[id] = tx
	return nil
}

// popCommand removes a <command> transaction from the map of in-flight commands.
func (c *client) popCommand(id string) (transaction, bool) {
	c.mCommands.Lock()
	defer c.mCommands.Unlock()
	tx, ok := c.commands[id]
	if ok {
		delete(c.commands, id)
	}
	return tx, ok
}

// cleanup cleans up and responds to all in-flight <hello> and <command> transactions.
// Each transaction will be finalized with err, which may be nil.
//
// TODO: clean up stale or abandoned transactions on a regular basis?
func (c *client) cleanup(err error) {
	c.mHellos.Lock()
	hellos := c.hellos
	c.hellos = nil
	c.mHellos.Unlock()
	for _, tx := range hellos {
		c.replyTo(tx, nil, err)
	}

	c.mCommands.Lock()
	commands := c.commands
	c.commands = nil
	c.mCommands.Unlock()
	for _, tx := range commands {
		c.replyTo(tx, nil, err)
	}
}
