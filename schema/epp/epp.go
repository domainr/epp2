package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// EPP represents an <epp> element as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP interface {
	Body() Body
}

type epp struct {
	B Body
}

func New(body Body) EPP {
	return &epp{body}
}

func (e *epp) Body() Body {
	return e.B
}

func (e *epp) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	start.Name.Space = NS
	start.Name.Local = "epp"
	type T epp
	return enc.EncodeElement((*T)(e), start)
}

func (e *epp) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.UseFactory(d, Schema, func(d *xml.Decoder) error {
		return schema.DecodeElements(d, func(v interface{}) error {
			if body, ok := v.(Body); ok {
				e.B = body
			}
			return nil
		})
	})
}
