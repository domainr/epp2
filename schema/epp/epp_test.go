package epp_test

import (
	"testing"

	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/schematest"
)

func TestEPPRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    string
		wantErr bool
	}{
		{
			`nil`,
			nil,
			``,
			false,
		},
		{
			`empty <epp> element`,
			&epp.EPP{},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"></epp>`,
			false,
		},
		{
			`<epp> with <hello> element`,
			&epp.EPP{Body: &epp.Hello{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello/></epp>`,
			false,
		},
		{
			`empty <greeting>`,
			&epp.EPP{Body: &epp.Greeting{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting></greeting></epp>`,
			false,
		},
		{
			`empty <command>`,
			&epp.EPP{Body: &epp.Command{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, nil, tt.v, tt.want, tt.wantErr)
		})
	}
}
