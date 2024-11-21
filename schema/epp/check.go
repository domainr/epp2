package epp

import (
	"github.com/domainr/epp2/internal/xml"

	"github.com/domainr/epp2/schema"
)

// Check represents an EPP <check> command as defined in [RFC 5730].
//
// [RFC 5730]: https://datatracker.ietf.org/doc/html/rfc5730#section-2.9.2.1
type Check struct {
	// XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	Check CheckType
}

func (Check) EPPAction() string { return "check" }

// MarshalXML implements the [xml.Marshaler] interface.
func (c *Check) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T Check
	return e.EncodeElement((*T)(c), schema.Rename(start, NS, c.EPPAction()))
}

// UnmarshalXML implements the [xml.Unmarshaler] interface. It requires an
// [xml.Decoder] with an associated [schema.Resolver] to correctly decode EPP
// <check> child elements.
func (c *Check) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.DecodeElements(d, func(v any) error {
		if check, ok := v.(CheckType); ok {
			c.Check = check
		}
		return nil
	})
}
