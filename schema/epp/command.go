package epp

import (
	"github.com/domainr/epp2/schema"
	"github.com/nbio/xml"
)

// Command represents an EPP client <command> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.5.
type Command struct {
	XMLName             struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`
	Command             CommandType
	ClientTransactionID string `xml:"clTRID,omitempty"`
}

func (Command) eppBody() {}

// UnmarshalXML implements the xml.Unmarshaler interface.
// It maps known EPP commands to their corresponding Go type.
func (c *Command) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Command
	var v struct {
		*T
		Command commandWrapper `xml:",any"`
	}
	v.T = (*T)(c)
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	c.Command = v.Command.Command
	return nil
}

// CommandType is a child element of EPP <Command>.
// Concrete CommandType types implement this interface.
type CommandType interface {
	eppCommand()
}

type commandWrapper struct {
	Command CommandType
}

func (c *commandWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.WithFactory(d, commandTypes, func(d *xml.Decoder) error {
		e, err := schema.DecodeElement(d, &start)
		if ct, ok := e.(CommandType); ok {
			c.Command = ct
		}
		return err
	})
}

var commandTypes = schema.FactoryFunc(func(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	case "check":
		return &Check{}
	case "create":
		return &Create{}
	case "delete":
		return &Delete{}
	case "info":
		return &Info{}
	case "login":
		return &Login{}
	case "logout":
		return &Logout{}
	case "poll":
		return &Poll{}
	case "renew":
		return &Renew{}
	case "transfer":
		return &Transfer{}
	case "update":
		return &Update{}
	}
	return nil
})
