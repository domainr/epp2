package schema

import (
	"reflect"
	"testing"

	"github.com/domainr/epp2/internal/xml"
)

func TestFactory(t *testing.T) {
	var n int
	f := FactoryFunc(func(name xml.Name) interface{} {
		if name.Space != "space" {
			return nil
		}
		switch name.Local {
		case "bytes":
			return []byte{}
		case "struct":
			return &struct{}{}
		case "int":
			var v int
			return &v
		}
		return nil
	})

	tests := []struct {
		name string
		arg  xml.Name
		want interface{}
	}{
		{
			`empty name`,
			xml.Name{},
			nil,
		},
		{
			`no namespace`,
			xml.Name{Local: "bytes"},
			nil,
		},
		{
			`bytes`,
			xml.Name{Space: "space", Local: "bytes"},
			[]byte{},
		},
		{
			`struct`,
			xml.Name{Space: "space", Local: "struct"},
			&struct{}{},
		},
		{
			`int`,
			xml.Name{Space: "space", Local: "int"},
			&n,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := f.New(tt.arg)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("f.New(%v)\nGot:  %#v\nWant: %#v", tt.arg, got, tt.want)
			}
		})
	}
}

type testFactory struct {
	v interface{}
}

func (f *testFactory) New(xml.Name) interface{} {
	return f.v
}

func TestFactories(t *testing.T) {
	a := &testFactory{}
	b := &testFactory{&struct{}{}}
	c := &testFactory{[]byte{}}

	tests := []struct {
		name string
		args []Factory
		want factories
	}{
		{
			`nil`,
			nil,
			nil,
		},
		{
			`empty slice`,
			factories{},
			factories{},
		},
		{
			`one element`,
			factories{a},
			factories{a},
		},
		{
			`two elements`,
			factories{a, b},
			factories{a, b},
		},
		{
			`three elements`,
			factories{a, b, c},
			factories{a, b, c},
		},
		{
			`mixed nils`,
			factories{a, nil, nil, b, c, nil},
			factories{a, b, c},
		},
		{
			`nested`,
			factories{factories{a, b}, c},
			factories{a, b, c},
		},
		{
			`nested with nils`,
			factories{factories{nil, a, b}, nil, nil, c, nil},
			factories{a, b, c},
		},
		{
			`deeply nested`,
			factories{factories{factories{factories{a}, b}}, c},
			factories{a, b, c},
		},
		{
			`deeply nested with nils`,
			factories{nil, factories{factories{nil, factories{a, nil, nil}, b}}, factories{}, nil, c, nil},
			factories{a, b, c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Factories(tt.args...)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Factories()\nGot:  %#v\nWant: %#v", got, tt.want)
			}
		})
	}
}
