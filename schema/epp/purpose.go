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

func parseOnePurpose(s string) Purpose {
	switch s {
	case "admin":
		return PurposeAdmin
	case "contact":
		return PurposeContact
	case "provisioning":
		return PurposeProvisioning
	case "other":
		return PurposeOther
	}
	return 0
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
func (p *Purpose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return schema.DecodeElements(d, func(v interface{}) error {
		if e, ok := v.(*schema.Any); ok && e.XMLName.Space == NS {
			*p |= parseOnePurpose(e.XMLName.Local)
		}
		return nil
	})
}
