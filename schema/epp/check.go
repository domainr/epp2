package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Check represents an EPP <check> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.1.
type Check struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	Check   CheckType
}

func (Check) eppCommand() {}

// UnmarshalXML implements the xml.Unmarshaler interface. It requires an
// xml.Decoder with an associated schema.Factory to correctly decode EPP <check>
// sub-elements.
func (c *Check) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	elements, err := schema.DecodeElements(d)
	if len(elements) > 0 {
		if check, ok := elements[0].(CheckType); ok {
			c.Check = check
		}
	}
	return err

}

// CheckType is a child element of EPP <check>.
// Concrete CheckType types implement this interface.
type CheckType interface {
	EPPCheck()
}
