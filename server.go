package epp

import (
	"context"
	"net"
	"sync"
	"sync/atomic"

	"github.com/domainr/epp2/protocol"
	"github.com/domainr/epp2/schema/epp"
)

// Server is an EPP version 1.0 server.
type Server struct {
	// Name is the name of this EPP server. It is sent to clients in a EPP
	// <greeting> message. If empty, a reasonable default will be used.
	Name string

	// Config describes the EPP server configuration. Configuration
	// parameters are announced to EPP clients in an EPP <greeting> message.
	Config Config

	// Handler is called in a goroutine for each incoming EPP connection.
	// The connection will be closed when Handler returns.
	Handler func(Session) error

	// ConnContext is called for each incoming connection.
	// If nil, context.Background() will be called.
	ConnContext func() context.Context

	inShutdown atomic.Bool

	mu            sync.Mutex
	listeners     map[net.Listener]struct{}
	listenerGroup sync.WaitGroup
}

func (s *Server) connContext() context.Context {
	if s.ConnContext != nil {
		return s.ConnContext()
	}
	return context.Background()
}

func (s *Server) shuttingDown() bool {
	return s.inShutdown.Load()
}

func (s *Server) trackListener(l net.Listener, add bool) bool {
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
func (s *Server) Serve(l net.Listener) error {
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
		go s.Accept(&protocol.NetConn{Conn: conn})
	}
}

// Accept accepts a connection and receives and processes EPP commands.
func (s *Server) Accept(conn protocol.Conn) error {
	if s.shuttingDown() {
		return ErrServerClosed
	}
	session := &session{
		ctx:  s.connContext(),
		conn: conn,
	}
	return s.handle(session)
}

func (s *Server) handle(sess *session) error {
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
	ctx  context.Context
	conn protocol.Conn
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
	return s.conn.Close()
}
