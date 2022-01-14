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

func (l *Login) ScanElement(start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
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

func (o *Outer) ScanElement(start xml.StartElement) (interface{}, error) {
	fmt.Println(start.Name.Local)
	switch start.Name.Local {
	case "inner":
		return &o.inner, nil
	}
	return nil, nil
}

func TestScan(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		want    interface{}
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := xml.NewDecoder(strings.NewReader(tt.xml))

			var got interface{}
			if tt.want != nil {
				got = reflect.New(reflect.TypeOf(tt.want).Elem()).Interface()
			}

			err := Scan(d, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Scan()\nGot:  %#v\nWant: %#v", got, tt.want)
			}
		})
	}
}

func TestScanFor(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		space   string
		local   string
		want    interface{}
		wantErr bool
	}{
		{
			`empty login`,
			`<login></login>`,
			"", "login",
			&Login{},
			false,
		},
		{
			`login with empty child tags`,
			`<login><clID></clID><pw></pw></login>`,
			"", "login",
			&Login{},
			false,
		},
		{
			`empty outer`,
			`<outer></outer>`,
			"", "outer",
			&Outer{},
			false,
		},
		{
			`outer with inner`,
			`<outer><inner></inner></outer>`,
			"", "outer",
			&Outer{},
			false,
		},
		{
			`outer with inner with value`,
			`<outer><inner>hello world</inner></outer>`,
			"", "outer",
			&Outer{"hello world"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := xml.NewDecoder(strings.NewReader(tt.xml))

			var got interface{}
			if tt.want != nil {
				got = reflect.New(reflect.TypeOf(tt.want).Elem()).Interface()
			}

			err := ScanFor(d, xml.Name{Space: tt.space, Local: tt.local}, got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Scan()\nGot:  %#v\nWant: %#v", got, tt.want)
			}
		})
	}
}
