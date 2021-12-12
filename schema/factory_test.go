package schema

import (
	"reflect"
	"testing"

	"github.com/nbio/xml"
)

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
