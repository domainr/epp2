package epp_test

import (
	"testing"

	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/schematest"
	"github.com/domainr/epp2/schema/std"
)

func TestResponseRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       any
		want    string
		wantErr bool
	}{
		{
			`empty <response>`,
			&epp.EPP{Body: &epp.Response{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`simple code 1000`,
			&epp.EPP{
				Body: &epp.Response{
					Results: []epp.Result{
						{
							Code:    epp.Success,
							Message: epp.Success.Message(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="1000"><msg lang="en">Command completed successfully</msg></result><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`multiple result codes`,
			&epp.EPP{
				Body: &epp.Response{
					Results: []epp.Result{
						{
							Code:    epp.ErrParameterRange,
							Message: epp.ErrParameterRange.Message(),
						},
						{
							Code:    epp.ErrParameterSyntax,
							Message: epp.ErrParameterSyntax.Message(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="2004"><msg lang="en">Parameter value range error</msg></result><result code="2005"><msg lang="en">Parameter value syntax error</msg></result><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`with extValue>reason`,
			&epp.EPP{
				Body: &epp.Response{
					Results: []epp.Result{
						{
							Code:    epp.ErrBillingFailure,
							Message: epp.ErrBillingFailure.Message(),
							ExtensionValues: []epp.ExtensionValue{
								{
									Reason: epp.Message{Lang: "en", Value: "Command exceeds available balance"},
								},
							},
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="2104"><msg lang="en">Billing failure</msg><extValue><reason lang="en">Command exceeds available balance</reason></extValue></result><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`with transaction IDs`,
			&epp.EPP{
				Body: &epp.Response{
					Results: []epp.Result{
						{
							Code:    epp.Success,
							Message: epp.Success.Message(),
						},
					},
					TransactionID: epp.TransactionID{
						Client: "12345",
						Server: "abcde",
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="1000"><msg lang="en">Command completed successfully</msg></result><trID><clTRID>12345</clTRID><svTRID>abcde</svTRID></trID></response></epp>`,
			false,
		},
		{
			`with basic <msgQ>`,
			&epp.EPP{
				Body: &epp.Response{
					MessageQueue: &epp.MessageQueue{Count: 5, ID: "67890"},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><msgQ count="5" id="67890"/><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`with <msgQ> with date`,
			&epp.EPP{
				Body: &epp.Response{
					MessageQueue: &epp.MessageQueue{
						Count: 5,
						ID:    "67890",
						Date:  std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><msgQ count="5" id="67890"><qDate>2000-01-01T00:00:00Z</qDate></msgQ><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`with full <msgQ>`,
			&epp.EPP{
				Body: &epp.Response{
					MessageQueue: &epp.MessageQueue{
						Count: 5,
						ID:    "67890",
						Date:  std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><msgQ count="5" id="67890"><qDate>2000-01-01T00:00:00Z</qDate></msgQ><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, nil, tt.v, tt.want, tt.wantErr)
		})
	}
}
