package epp

import (
	"github.com/domainr/epp2/schema/std"
	"github.com/nbio/xml"
)

// Greeting represents an EPP server <greeting> message as defined in RFC 5730.
type Greeting struct {
	XMLName     struct{}     `xml:"urn:ietf:params:xml:ns:epp-1.0 greeting"`
	ServerName  string       `xml:"svID,omitempty"`
	ServerDate  *std.Time    `xml:"svDate"`
	ServiceMenu *ServiceMenu `xml:"svcMenu"`
	DCP         *DCP         `xml:"dcp"`
}

func (Greeting) eppBody() {}

// ServiceMenu represents an EPP <svcMenu> element as defined in RFC 5730.
type ServiceMenu struct {
	Versions         []string          `xml:"version"`
	Languages        []string          `xml:"lang"`
	Objects          []string          `xml:"objURI"`
	ServiceExtension *ServiceExtension `xml:"svcExtension"`
}

// DCP represents a server data collection policy as defined in RFC 5730.
type DCP struct {
	Access     Access      `xml:"access"`
	Statements []Statement `xml:"statement"`
	Expiry     *Expiry     `xml:"expiry"`
}

// Statement describes an EPP server’s data collection purpose, receipient(s), and retention policy.
type Statement struct {
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
}

// Purpose represents an EPP server’s purpose for data collection.
type Purpose struct {
	Admin        std.Bool `xml:"admin"`
	Contact      std.Bool `xml:"contact"`
	Provisioning std.Bool `xml:"provisioning"`
	Other        std.Bool `xml:"other"`
}

func PurposeNone() Purpose         { return Purpose{} }
func PurposeAdmin() Purpose        { return Purpose{Admin: true} }
func PurposeContact() Purpose      { return Purpose{Contact: true} }
func PurposeProvisioning() Purpose { return Purpose{Provisioning: true} }
func PurposeOther() Purpose        { return Purpose{Other: true} }
func PurposeAll() Purpose          { return Purpose{true, true, true, true} }

// Recipient represents an EPP server’s purpose for data collection.
type Recipient struct {
	Other     std.Bool `xml:"other"`
	Ours      *Ours    `xml:"ours"`
	Public    std.Bool `xml:"public"`
	Same      std.Bool `xml:"same"`
	Unrelated std.Bool `xml:"unrelated"`
}

func RecipientNone() Recipient                 { return Recipient{} }
func RecipientOther() Recipient                { return Recipient{Other: true} }
func RecipientOurs(recipient string) Recipient { return Recipient{Ours: &Ours{recipient}} }
func RecipientPublic() Recipient               { return Recipient{Public: true} }
func RecipientSame() Recipient                 { return Recipient{Same: true} }
func RecipientUnrelated() Recipient            { return Recipient{Unrelated: true} }

// Ours represents an EPP server’s description of an <ours> recipient.
type Ours struct {
	Recipient string `xml:"recDesc"`
}

// MarshalXML impements the xml.Marshaler interface.
// Writes a single self-closing <ours/> if v.Recipient is not set.
func (v *Ours) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if v.Recipient == "" {
		return e.EncodeToken(xml.SelfClosingElement(start))
	}
	type T Ours
	return e.EncodeElement((*T)(v), start)
}

// Expiry defines an EPP server’s data retention duration.
type Expiry struct {
	Absolute *std.Time     `xml:"absolute"`
	Relative *std.Duration `xml:"relative"`
}
