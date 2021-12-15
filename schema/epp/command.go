package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// Command represents an EPP client <command> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.5.
type Command struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`

	// Action is an element whose tag corresponds to one of the valid EPP
	// commands described in RFC 5730. The command element MAY contain
	// either protocol-specified or object-specified child elements.
	Action Action

	// Extensions is an OPTIONAL <extension> element that MAY be used for
	// server- defined command extensions.
	Extensions []Extension `xml:"extension,omitempty"`

	// ClientTransactionID is an OPTIONAL <clTRID> (client transaction
	// identifier) element that MAY be used to uniquely identify the command
	// to the client. Clients are responsible for maintaining their own
	// transaction identifier space to ensure uniqueness.
	ClientTransactionID string `xml:"clTRID,omitempty"`
}

func (Command) eppBody() {}

// UnmarshalXML implements the xml.Unmarshaler interface.
// It maps known EPP commands to their corresponding Go type.
// It requires an xml.Decoder with an associated schema.Factory to
// correctly decode EPP <command> sub-elements.
func (c *Command) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Command
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

// UnmarshalXML requires an xml.Decoder with an associated schema.Factory to
// property decode EPP <command> actions.
func (w *actionWrapper) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	v, err := schema.DecodeElement(d, start)
	if a, ok := v.(Action); ok {
		w.Action = a
	}
	return err
}
