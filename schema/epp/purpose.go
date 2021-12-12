package epp

import (
	"strings"

	"github.com/domainr/epp2/schema"
	"github.com/nbio/xml"
)

// Purpose represents an EPP server’s purpose for data collection.
// Multiple values of Purpose can be or’ed together.
type Purpose uint8

const (
	PurposeAdmin        Purpose = 1
	PurposeContact      Purpose = 2
	PurposeProvisioning Purpose = 4
	PurposeOther        Purpose = 8
)

// ParsePurpose parses s into an Purpose.
// It returns PurposeOther if s is not recognized.
func ParsePurpose(s string) Purpose {
	var p Purpose
	for _, t := range strings.Split(s, " ") {
		switch t {
		case "admin":
			p |= PurposeAdmin
		case "contact":
			p |= PurposeContact
		case "provisioning":
			p |= PurposeProvisioning
		case "other":
			p |= PurposeOther
		}
	}
	if p == 0 {
		p = PurposeOther
	}
	return p
}

// String returns the a string representation for p. If p has only one value,
// the returned string can be used as an XML tag name.
func (p Purpose) String() string {
	var a [4]string
	s := a[:0]
	if p&PurposeAdmin != 0 {
		s = append(s, "admin")
	}
	if p&PurposeContact != 0 {
		s = append(s, "contact")
	}
	if p&PurposeProvisioning != 0 {
		s = append(s, "provisioning")
	}
	if p&PurposeOther != 0 {
		s = append(s, "other")
	}
	return strings.Join(s, " ")
}

// MarshalXML impements the xml.Marshaler interface.
func (p *Purpose) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Space = NS
	start.Name.Local = "purpose"
	err := e.EncodeToken(start)
	if err != nil {
		return nil
	}
	for i := PurposeAdmin; i <= PurposeOther; i <<= 1 {
		if *p&i != 0 {
			err = e.EncodeToken(xml.SelfClosingElement{Name: xml.Name{Space: NS, Local: i.String()}})
			if err != nil {
				return nil
			}
		}
	}
	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

// UnmarshalXML implements the xml.Unmarshaler interface.
func (a *Purpose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.DecodeElements(d, func(v interface{}) error {
		if e, ok := v.(*schema.Any); ok && e.XMLName.Space == NS {
			*a |= ParsePurpose(e.XMLName.Local)
		}
		return nil
	})
}
