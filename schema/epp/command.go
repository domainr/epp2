package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// Command represents an EPP client <command> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.5.
type Command struct {
	// Action is an element whose tag corresponds to one of the valid EPP
	// commands described in RFC 5730. The command element MAY contain
	// either protocol-specified or object-specified child elements.
	Action Action

	// Extensions is an OPTIONAL <extension> element that MAY be used for
	// server- defined command extensions.
	Extensions Extensions

	// ClientTransactionID is an OPTIONAL <clTRID> (client transaction
	// identifier) element that MAY be used to uniquely identify the command
	// to the client. Clients are responsible for maintaining their own
	// transaction identifier space to ensure uniqueness.
	ClientTransactionID string
}

func (Command) eppBody() {}

type commandXML struct {
	Action              Action
	Extensions          Extensions `xml:"extension,omitempty"`
	ClientTransactionID string     `xml:"clTRID,omitempty"`
}

// MarshalXML implements the [xml.Marshaler] interface.
func (c *Command) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.EncodeElement((*commandXML)(c), schema.Rename(start, NS, "command"))
}

// UnmarshalXML implements the xml.Unmarshaler interface.
// It maps known EPP commands to their corresponding Go type.
// It requires an xml.Decoder with an associated schema.Resolver to
// correctly decode EPP <command> sub-elements.
func (c *Command) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T commandXML
	var v struct {
		*T
		V actionWrapper `xml:",any"`
	}
	v.T = (*T)(c)
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	c.Action = v.V.Action
	return nil
}

type actionWrapper struct {
	Action Action
}

// UnmarshalXML requires an xml.Decoder with an associated schema.Resolver to
// property decode EPP <command> actions.
func (w *actionWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	v, err := schema.DecodeElement(d, start)
	if a, ok := v.(Action); ok {
		w.Action = a
	}
	return err
}
