package schema

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/domainr/epp2/internal/xml"
)

type Login struct {
	user    string
	pass    string
	newPass *string
}

func (l *Login) ScanStartElement(start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "login":
		return l, nil
	case "clID":
		return &l.user, nil
	case "pw":
		return &l.pass, nil
	case "newPW":
		l.newPass = new(string)
		return l.newPass, nil
	}
	return nil, nil
}

type Outer struct {
	inner string
}

func (o *Outer) ScanStartElement(start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "outer":
		return o, nil
	case "inner":
		return &o.inner, nil
	}
	return nil, nil
}

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		v       interface{}
		want    interface{}
		wantErr bool
	}{
		{
			`nil`,
			``,
			nil,
			nil,
			false,
		},
		{
			`unbalanced end tag`,
			`</a>`,
			nil,
			nil,
			true,
		},
		{
			`incorrect end tag`,
			`<a></b>`,
			nil,
			nil,
			true,
		},
		{
			`empty login`,
			`<login></login>`,
			&Login{},
			&Login{},
			false,
		},
		{
			`login with empty child tags`,
			`<login><clID></clID><pw></pw></login>`,
			&Login{},
			&Login{},
			false,
		},
		{
			`empty outer`,
			`<outer></outer>`,
			&Outer{},
			&Outer{},
			false,
		},
		{
			`outer with inner`,
			`<outer><inner></inner></outer>`,
			&Outer{},
			&Outer{},
			false,
		},
		{
			`outer with inner with value`,
			`<outer><inner>hello world</inner></outer>`,
			&Outer{},
			&Outer{"hello world"},
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

			if !reflect.DeepEqual(tt.want, tt.v) {
				t.Errorf("Scan()\nGot:  %#v\nWant: %#v", tt.v, tt.want)
			}
		})
	}
}
