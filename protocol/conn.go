package protocol

import (
	"encoding/xml"
	"sync"

	"github.com/domainr/epp2/protocol/wire"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// Conn represents a low-level EPP connection.
// Reads and writes are synchronized, so a Conn is safe to use from multiple goroutines.
type Conn interface {
	// ReadEPP reads the next EPP message from the underlying connection.
	// An error will be returned if the underlying connection is closed or an error occurs
	// reading from the connection.
	ReadEPP() (epp.Body, error)

	// WriteEPP writes an EPP response to the underlying connection. An error will
	// be returned if the underlying connection is closed or an error occurs
	// writing to the connection.
	WriteEPP(epp.Body) error

	// Close closes the connection.
	// No attempt is made to wait for or clean up any transactions in flight.
	Close() error
}

type eppConn struct {
	// reading synchronizes reads from conn.
	reading sync.Mutex

	// writing synchronizes writes to conn.
	writing sync.Mutex

	// conn holds the underlying data unit connection.
	conn wire.Conn

	schemas schema.Schemas
}

var _ Conn = &eppConn{}

// NewConn returns a new [Conn] using conn as the underlying transport.
//
// Messages from the peer will be decoded using [schemas.Schema] schemas.
// If no schemas are provided, a set of reasonable defaults will be used.
func NewConn(conn wire.Conn, schemas schema.Schemas) *eppConn {
	return &eppConn{
		conn:    conn,
		schemas: schemas,
	}
}

func (c *eppConn) ReadEPP() (epp.Body, error) {
	c.reading.Lock()
	data, err := c.conn.ReadDataUnit()
	c.reading.Unlock()
	if err != nil {
		return nil, err
	}
	var e epp.EPP
	err = schema.Unmarshal(data, &e, c.schemas)
	return e.Body, err
}

func (c *eppConn) WriteEPP(body epp.Body) error {
	e := epp.EPP{Body: body}
	x, err := xml.Marshal(&e) // TODO: implement schema.Marshal()
	if err != nil {
		return err
	}
	c.writing.Lock()
	err = c.conn.WriteDataUnit(x)
	c.writing.Unlock()
	return err
}

func (c *eppConn) Close() error {
	return c.conn.Close()
}
