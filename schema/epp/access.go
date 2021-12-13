package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Access represents an EPP server’s scope of data access as defined in RFC 5730.
type Access string

const (
	AccessNull             Access = "null"
	AccessNone             Access = "none"
	AccessPersonal         Access = "personal"
	AccessOther            Access = "other"
	AccessPersonalAndOther Access = "personalAndOther"
	AccessAll              Access = "all"
)

func parseAccess(s string) Access {
	switch s {
	case "null", "none", "personal", "personalAndOther", "all":
		return Access(s)
	}
	return ""
}

// MarshalXML impements the xml.Marshaler interface.
func (a Access) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	err := e.EncodeToken(start)
	if err != nil {
		return nil
	}
	if parseAccess(string(a)) != "" {
		err = e.EncodeToken(xml.SelfClosingElement{Name: xml.Name{Space: NS, Local: string(a)}})
		if err != nil {
			return nil
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (a *Access) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.DecodeElements(d, func(v interface{}) error {
		if e, ok := v.(*schema.Any); ok && e.XMLName.Space == NS {
			*a = parseAccess(e.XMLName.Local)
		}
		return nil
	})
}
