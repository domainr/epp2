package epp_test

import (
	"testing"

	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/schematest"
)

func TestEPPExtensionsRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    string
		wantErr bool
	}{
		{
			`<epp> with empty <extension> element`,
			&epp.EPP{Body: &epp.Extensions{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><extension></extension></epp>`,
			false,
		},
		{
			`<epp> with one <extension> sub-element`,
			&epp.EPP{Body: &epp.Extensions{&fooBar{}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><extension><foo:bar xmlns:foo="urn:example:foo-1.0"></foo:bar></extension></epp>`,
			false,
		},
		{
			`<epp> with two <extension> sub-elements`,
			&epp.EPP{Body: &epp.Extensions{&fooBar{}, &fooBaz{}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><extension><foo:bar xmlns:foo="urn:example:foo-1.0"></foo:bar><foo:baz xmlns:foo="urn:example:foo-1.0"></foo:baz></extension></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, fooSchema, tt.v, tt.want, tt.wantErr)
		})
	}
}

const fooNS = "urn:example:foo-1.0"

const fooSchema fooSchemaString = "foo"

var _ schema.Schema = fooSchema

type fooSchemaString string

func (o fooSchemaString) SchemaName() string {
	return string(o)
}

func (fooSchemaString) SchemaNS() []string {
	return []string{fooNS}
}

func (fooSchemaString) ResolveXML(name xml.Name) any {
	if name.Space != fooNS {
		return nil
	}
	switch name.Local {
	case "bar":
		return &fooBar{}
	case "baz":
		return &fooBaz{}
	}
	return nil
}

type fooBar struct {
	XMLName struct{} `xml:"urn:example:foo-1.0 foo:bar"`
}

func (fooBar) EPPExtension() {}

type fooBaz struct {
	XMLName struct{} `xml:"urn:example:foo-1.0 foo:baz"`
}

func (fooBaz) EPPExtension() {}
