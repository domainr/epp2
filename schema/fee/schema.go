//go:build ignore

package fee

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

// NS defines the IETF URN for the EPP fee 1.0 namespace.
// See https://www.iana.org/assignments/xml-registry/ns/epp/fee-1.0.txt.
const NS = "urn:ietf:params:xml:ns:epp:fee-1.0"

// Schema implements the schema.Schema interface for the EPP common namespace.
const Schema schemaString = "fee"

var _ schema.Schema = Schema

type schemaString string

func (o schemaString) SchemaName() string {
	return string(o)
}

func (schemaString) SchemaNS() []string {
	return []string{NS}
}

func (schemaString) ResolveXML(name xml.Name) any {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	// TODO: what are EPP fee types?
	}
	return nil
}

type Check struct{}

func (Check) EPPExtension() {}

// MarshalXML implements the [xml.Marshaler] interface.
func (c *Check) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T Check
	return e.EncodeElement((*T)(c), schema.Rename(start, NS, string(Schema)+":check"))
}

type CheckData struct{}

func (CheckData) EPPExtension() {}

type Create Transform[epp.Create]

type CreateData struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp:fee-1.0 fee:creData"`
	TransformData[any]
}

type Renew Transform[epp.Renew]

type RenewData struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp:fee-1.0 fee:renData"`
	TransformData[any]
}

type Transfer Transform[epp.Transfer]

type TransferData struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp:fee-1.0 fee:trnData"`
	TransformData[any]
}

type Update Transform[epp.Update]

type UpdateData struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp:fee-1.0 fee:updData"`
	TransformData[any]
}

type Transform[A epp.Action] struct{}

func (Transform[A]) EPPExtension() {}

// MarshalXML implements the [xml.Marshaler] interface.
func (t *Transform[A]) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T Transform[A]
	var a A
	return e.EncodeElement((*T)(t), schema.Rename(start, NS, a.EPPAction()))
}

type TransformData[A epp.ResponseData] struct{}

func (TransformData[A]) EPPExtension() {}
