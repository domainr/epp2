package epp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema/epp"
)

// server is an EPP version 1.0 server.
type server struct {
	// Name is the name of this EPP server. It is sent to clients in a EPP
	// <greeting> message. If empty, a reasonable default will be used.
	Name string

	// Config describes the EPP server configuration. Configuration
	// parameters are announced to EPP clients in an EPP <greeting> message.
	Config Config

	// Handler is called in a goroutine for each incoming EPP connection.
	// The connection will be closed when Handler returns.
	Handler func(Session) error

	inShutdown atomic.Bool

	mu            sync.Mutex
	listeners     map[net.Listener]struct{}
	listenerGroup sync.WaitGroup
}

func (s *server) shuttingDown() bool {
	return s.inShutdown.Load()
}

func (s *server) trackListener(l net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[net.Listener]struct{})
	}
	if add {
		if s.shuttingDown() {
			return false
		}
		s.listeners[l] = struct{}{}
		s.listenerGroup.Add(1)
	} else {
		delete(s.listeners, l)
		s.listenerGroup.Done()
	}
	return true
}

// Serve accepts incoming connections on [net.Listener] l,
// creating a new service goroutine for each connection.
// The service goroutines read commands and then call s.Handler to reply to them.
func (s *server) Serve(l net.Listener) error {
	if !s.trackListener(l, true) {
		return ErrServerClosed
	}
	defer s.trackListener(l, false)
	for {
		conn, err := l.Accept()
		if err != nil {
			if s.shuttingDown() {
				return ErrServerClosed
			}
			return err
		}
		var _ = conn
		// pconn := protocol.NewConn(&dataunit.NetConn{Conn: conn}, s.Config.Schemas)
		// go s.Handle(pconn)
	}
}

// Handle accepts a connection and receives and processes EPP commands.
func (s *server) Handle(conn net.Conn) error {
	if s.shuttingDown() {
		return ErrServerClosed
	}
	session := &session{
		// ctx:  s.connContext(conn),
		// conn: conn,
	}
	return s.handle(session)
}

func (s *server) handle(sess *session) error {
	defer sess.Close()
	if s.Handler == nil {
		return nil
	}
	return s.Handler(sess)
}

type Session interface {
	// Context returns the connection Context for this session. The Context
	// will be canceled if the underlying connection goes away or is closed.
	Context() context.Context

	// ReadCommand reads the next EPP command from the client. An error will
	// be returned if the underlying connection is closed or an error occurs
	// reading from the connection.
	ReadCommand() (*epp.Command, error)

	// WriteResponse sends an EPP response to the client. An error will
	// be returned if the underlying connection is closed or an error occurs
	// writing to the connection.
	WriteResponse(*epp.Response) error

	// Close closes the session and the underlying connection.
	Close() error
}

type session struct {
	ctx context.Context
	s   dataunit.Server // FIXME: this should be a protocol.Server or protocol.Session
}

var _ Session = &session{}

func (s *session) Context() context.Context {
	return s.ctx
}

func (s *session) ReadCommand() (*epp.Command, error) {
	// TODO: implement this
	return nil, nil
}

func (s *session) WriteResponse(r *epp.Response) error {
	// TODO: implement this
	return nil
}

func (s *session) Close() error {
	// return s.s.Close()
	return nil
}
