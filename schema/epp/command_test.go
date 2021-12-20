package epp_test

import (
	"testing"

	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/test"
)

func TestCommandRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`empty <command>`,
			epp.New(&epp.Command{}),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
		{
			`empty <check> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Check{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check></check></command></epp>`,
			false,
		},
		{
			`empty <create> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Create{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create></create></command></epp>`,
			false,
		},
		{
			`empty <delete> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Delete{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><delete></delete></command></epp>`,
			false,
		},
		{
			`empty <info> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Info{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><info></info></command></epp>`,
			false,
		},
		{
			`empty <login> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Login{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID></clID><pw></pw><options><version></version></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`empty <logout> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Logout{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><logout/></command></epp>`,
			false,
		},
		{
			`empty <poll> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Poll{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><poll></poll></command></epp>`,
			false,
		},
		{
			`empty <renew> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Renew{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><renew></renew></command></epp>`,
			false,
		},
		{
			`empty <transfer> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Transfer{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer></transfer></command></epp>`,
			false,
		},
		{
			`empty <update> command`,
			epp.New(
				&epp.Command{
					Action: &epp.Update{},
				},
			),
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><update></update></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, nil, tt.v, tt.want, tt.wantErr)
		})
	}
}
