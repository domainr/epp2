package epp_test

import (
	"testing"

	"github.com/domainr/epp2/schema/contact"
	"github.com/domainr/epp2/schema/domain"
	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/host"
	"github.com/domainr/epp2/schema/schematest"
	"github.com/domainr/epp2/schema/std"
)

func TestGreetingRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    string
		wantErr bool
	}{
		{
			`empty <greeting>`,
			&epp.EPP{Body: &epp.Greeting{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting></greeting></epp>`,
			false,
		},
		{
			`simple <greeting>`,
			&epp.EPP{
				Body: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate></greeting></epp>`,
			false,
		},
		{
			`complex <greeting>`,
			&epp.EPP{
				Body: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{contact.NS, domain.NS, host.NS},
					},
					DCP: &epp.DCP{},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI></svcMenu><dcp><access></access></dcp></greeting></epp>`,
			false,
		},
		{
			`complex <greeting> with complex <dcp>`,
			&epp.EPP{
				Body: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{contact.NS, domain.NS, host.NS},
					},
					DCP: &epp.DCP{
						Access: epp.AccessPersonalAndOther,
						Statements: []epp.Statement{
							{
								Purpose:   epp.PurposeAdmin(),
								Recipient: epp.Recipient{Ours: &epp.Ours{Recipient: "Domainr"}, Public: true},
							},
							{
								Purpose:   epp.Purpose{Contact: true, Other: true},
								Recipient: epp.Recipient{Other: true, Ours: &epp.Ours{}, Public: true},
							},
						},
						Expiry: &epp.Expiry{
							Relative: std.ParseDuration("P1Y").Pointer(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI></svcMenu><dcp><access><personalAndOther/></access><statement><purpose><admin/></purpose><recipient><ours><recDesc>Domainr</recDesc></ours><public/></recipient></statement><statement><purpose><contact/><other/></purpose><recipient><other/><ours/><public/></recipient></statement><expiry><relative>P365DT5H49M12S</relative></expiry></dcp></greeting></epp>`,
			false,
		},
		{
			`<greeting> with <dcp> with absolute expiry`,
			&epp.EPP{
				Body: &epp.Greeting{
					DCP: &epp.DCP{
						Expiry: &epp.Expiry{
							Absolute: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><dcp><access></access><expiry><absolute>2000-01-01T00:00:00Z</absolute></expiry></dcp></greeting></epp>`,
			false,
		},
		{
			`complex <greeting> with extensions`,
			&epp.EPP{
				Body: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{contact.NS, domain.NS, host.NS},
						ServiceExtension: &epp.ServiceExtension{
							Extensions: []string{
								"urn:ietf:params:xml:ns:fee-0.8",
								"urn:ietf:params:xml:ns:epp:fee-1.0",
							},
						},
					},
					DCP: &epp.DCP{
						Access: epp.AccessNull,
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI><svcExtension><extURI>urn:ietf:params:xml:ns:fee-0.8</extURI><extURI>urn:ietf:params:xml:ns:epp:fee-1.0</extURI></svcExtension></svcMenu><dcp><access><null/></access></dcp></greeting></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, nil, tt.v, tt.want, tt.wantErr)
		})
	}
}
