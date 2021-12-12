package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Access represents an EPP server’s scope of data access as defined in RFC 5730.
type Access uint8

const (
	AccessNull             Access = 0
	AccessNone             Access = 1
	AccessPersonal         Access = 2
	AccessOther            Access = 4
	AccessPersonalAndOther Access = AccessPersonal | AccessOther
	AccessAll              Access = 8
)

// ParseAccess parses s into an Access.
// It returns AccessNull if s is not recognized.
func ParseAccess(s string) Access {
	switch s {
	case "null":
		return AccessNull
	case "none":
		return AccessNone
	case "personal":
		return AccessPersonal
	case "personalAndOther":
		return AccessPersonalAndOther
	case "all":
		return AccessAll
	}
	return AccessNull
}

// String returns the EPP tag name for a.
func (a Access) String() string {
	switch a {
	case AccessNull:
		return "null"
	case AccessNone:
		return "none"
	case AccessPersonal:
		return "personal"
	case AccessPersonalAndOther:
		return "personalAndOther"
	case AccessAll:
		return "all"
	}
	return ""
}

// MarshalXML impements the xml.Marshaler interface.
func (a *Access) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Space = NS
	start.Name.Local = "access"
	err := e.EncodeToken(start)
	if err != nil {
		return nil
	}
	local := a.String()
	if local != "" {
		err = e.EncodeToken(xml.SelfClosingElement{Name: xml.Name{Space: NS, Local: local}})
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
			*a = ParseAccess(e.XMLName.Local)
		}
		return nil
	})
}
