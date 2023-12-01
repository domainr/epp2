package epp

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
)

// Command represents an EPP <extension> element as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.7.1.
type Extensions struct {
	XMLName    struct{}    `xml:"extension"`
	Extensions []Extension `xml:",omitempty"`
}

func (*Extensions) eppBody() {}

// UnmarshalXML implements the xml.Unmarshaler interface. It requires an
// xml.Decoder with an associated schema.Resolver to correctly decode EPP <extension>
// sub-elements.
func (exts *Extensions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.DecodeElements(d, func(v any) error {
		if ext, ok := v.(Extension); ok {
			exts.Extensions = append(exts.Extensions, ext)
		}
		return nil
	})
}
