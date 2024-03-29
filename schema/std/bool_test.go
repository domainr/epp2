package std

import (
	"testing"

	"github.com/domainr/epp2/schema/schematest"
)

func TestBool(t *testing.T) {
	type T1 struct {
		XMLName struct{} `xml:"example"`
		Fred    Bool     `xml:"fred"`
		Jane    Bool     `xml:"jane"`
		Susan   Bool     `xml:"susan"`
	}

	type T2 struct {
		XMLName struct{} `xml:"example,selfclosing"`
		Fred    Bool     `xml:"fred,attr"`
		Jane    Bool     `xml:"jane,attr,omitempty"`
		Susan   Bool     `xml:"susan,attr,omitempty"`
	}

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
			`no tags`,
			&T1{},
			`<example></example>`,
			false,
		},
		{
			`Fred`,
			&T1{Fred: true},
			`<example><fred/></example>`,
			false,
		},
		{
			`Jane`,
			&T1{Jane: true},
			`<example><jane/></example>`,
			false,
		},
		{
			`Fred and Susan`,
			&T1{Fred: true, Susan: true},
			`<example><fred/><susan/></example>`,
			false,
		},
		{
			`Fred attribute`,
			&T2{Fred: true},
			`<example fred="1"/>`,
			false,
		},
		{
			`Jane attribute`,
			&T2{Jane: true},
			`<example fred="0" jane="1"/>`,
			false,
		},
		{
			`Fred and Susan attributes`,
			&T2{Fred: true, Susan: true},
			`<example fred="1" susan="1"/>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			schematest.RoundTrip(t, nil, tt.v, tt.want, tt.wantErr)
		})
	}
}
