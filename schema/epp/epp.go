package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// EPP represents an <epp> element as defined in [RFC 5730].
//
// [RFC 5730]: https://datatracker.ietf.org/doc/rfc5730/
type EPP struct {
	// Body is any valid EPP child element.
	Body Body
}

// MarshalXML implements the [xml.Marshaler] interface.
func (e *EPP) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type T EPP
	return enc.EncodeElement((*T)(e), schema.Rename(start, NS, "epp"))
}

// UnmarshalXML implements the [xml.Unmarshaler] interface.
func (e *EPP) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.UseResolver(d, Schema, func(d *xml.Decoder) error {
		return schema.DecodeElements(d, func(v any) error {
			if body, ok := v.(Body); ok {
				e.Body = body
			}
			return nil
		})
	})
}
