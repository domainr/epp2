package epp_test

import (
	"testing"

	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/domain"
	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/schematest"
)

func TestCheckRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		resolver schema.Resolver
		v        any
		want     string
		wantErr  bool
	}{
		{
			`empty <domain:check> command`,
			domain.Schema,
			&epp.EPP{
				Body: &epp.Command{
					Action: &epp.Check{
						Check: &domain.Check{},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"></domain:check></check></command></epp>`,
			false,
		},
		{
			`single <domain:check> command`,
			domain.Schema,
			&epp.EPP{
				Body: &epp.Command{
					Action: &epp.Check{
						Check: &domain.Check{
							Names: []string{"example.com"},
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:check></check></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, tt.resolver, tt.v, tt.want, tt.wantErr)
		})
	}
}
