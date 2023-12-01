package protocol

import (
	"github.com/domainr/epp2/protocol/dataunit"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Writer is the interface implemented by any type that can write an EPP message body.
type Writer interface {
	WriteEPP(epp.Body) error
}

type writerFunc func(body epp.Body) error

func (f writerFunc) WriteEPP(body epp.Body) error {
	return f(body)
}

// Server is a low-level server for the Extensible Provisioning Protocol (EPP)
// as defined in [RFC 5730]. A Server is safe to use from multiple goroutines.
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type Server interface {
	// ServeEPP provides an client EPP request and a mechanism to respond to the request.
	// It blocks until a response is received or the underlying connection is closed.
	// The returned [Writer] should only be used once. The returned Writer will always
	// be non-nil, so the caller can respond to a malformed client request.
	ServeEPP() (epp.Body, Writer, error)

	// Close closes the connection.
	Close() error
}

type server struct {
	server dataunit.Server
	coder  coder
}

// Serve services conn as an EPP server, sending greeting as the initial <greeting>
// message to the client.
// EPP requests from the client will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func Serve(conn dataunit.Conn, greeting epp.Body, schemas ...schema.Schema) (Server, error) {
	s := newServer(conn, schemas)
	// Send the initial <greeting> to the client.
	data, err := s.coder.marshalXML(greeting)
	if err != nil {
		return nil, err
	}
	return s, conn.WriteDataUnit(data)
}

func newServer(conn dataunit.Conn, schemas schema.Schemas) *server {
	if len(schemas) == 0 {
		schemas = DefaultSchemas()
	}
	return &server{
		server: dataunit.Server{Conn: conn},
		coder:  coder{schemas},
	}
}

// Close closes the connection, interrupting any in-flight requests.
func (s *server) Close() error {
	return s.server.Conn.Close()
}

func (s *server) ServeEPP() (epp.Body, Writer, error) {
	data, w, err := s.server.ServeDataUnit()
	f := writerFunc(func(body epp.Body) error {
		data, err := s.coder.marshalXML(body)
		if err != nil {
			return err
		}
		return w.WriteDataUnit(data)
	})
	if err != nil {
		return nil, f, err
	}
	body, err := s.coder.umarshalXML(data)
	return body, f, err
}