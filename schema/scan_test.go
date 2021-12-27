package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/domainr/epp2/internal/xml"
)

type Login struct {
	Username    string
	Password    string
	NewPassword *string
}

func (l *Login) ScanStartElement(r xml.TokenReader, start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "login":
		return l, nil
	case "clID":
		return &l.Username, nil
	case "pw":
		return &l.Password, nil
	case "newPW":
		l.NewPassword = new(string)
		return l.NewPassword, nil
	}
	return nil, nil
}

type Outer struct {
	inner Inner
}

func (o *Outer) ScanStartElement(r xml.TokenReader, start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "outer":
		return o, nil
	case "inner":
		return &o.inner, nil
	}
	return nil, nil
}

type Inner struct {
	v string
}

type Outer2 Outer

func (o *Outer2) ScanStartElement(r xml.TokenReader, start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "outer":
		return o, nil
	case "inner":
		return nil, Scan(r, &o.inner)
	}
	return nil, nil
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
		{
			`empty outer (recursive)`,
			`<outer></outer>`,
			&Outer2{},
			false,
		},
		{
			`outer (recursive) with inner`,
			`<outer><inner></inner></outer>`,
			&Outer2{},
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
