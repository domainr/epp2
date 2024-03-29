package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// EPP represents an <epp> element as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Body is any valid EPP child element.
	Body Body
}

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
