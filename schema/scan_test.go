package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/domainr/epp2/internal/xml"
)

type Login struct {
	user    string
	pass    string
	newPass *string
}

func (l *Login) ScanStartElement(r xml.TokenReader, start xml.StartElement) error {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "login":
		return Scan(r, l)
	case "clID":
		return Scan(r, &l.user)
	case "pw":
		return Scan(r, &l.pass)
	case "newPW":
		l.newPass = new(string)
		return Scan(r, l.newPass)
	}
	return nil
}

type Outer struct {
	inner Inner
}

func (o *Outer) ScanStartElement(r xml.TokenReader, start xml.StartElement) error {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "outer":
		return Scan(r, o)
	case "inner":
		return Scan(r, &o.inner)
	}
	return nil
}

type Inner struct {
	v string
}

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		v       interface{}
		wantErr bool
	}{
		{
			`nil`,
			``,
			nil,
			false,
		},
		{
			`unbalanced end tag`,
			`</a>`,
			nil,
			true,
		},
		{
			`incorrect end tag`,
			`<a></b>`,
			nil,
			true,
		},
		{
			`empty login`,
			`<login></login>`,
			&Login{},
			false,
		},
		{
			`login with empty child tags`,
			`<login><clID></clID><pw></pw></login>`,
			&Login{},
			false,
		},
		{
			`empty outer`,
			`<outer></outer>`,
			&Outer{},
			false,
		},
		{
			`outer with inner`,
			`<outer><inner></inner></outer>`,
			&Outer{},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := xml.NewDecoder(strings.NewReader(tt.xml))
			err := Scan(d, tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
