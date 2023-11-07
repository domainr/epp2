package schema

import (
	"reflect"
	"testing"

	"github.com/domainr/epp2/internal/xml"
)

func TestResolver(t *testing.T) {
	var n int
	f := ResolverFunc(func(name xml.Name) any {
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
		want any
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
			got := f.ResolveXML(tt.arg)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("ResolveXML(%v)\nGot:  %#v\nWant: %#v", tt.arg, got, tt.want)
			}
		})
	}
}

type testResolver struct {
	v any
}

func (f *testResolver) ResolveXML(xml.Name) any {
	return f.v
}

func TestFlatten(t *testing.T) {
	a := &testResolver{}
	b := &testResolver{&struct{}{}}
	c := &testResolver{[]byte{}}

	tests := []struct {
		name string
		args []Resolver
		want resolvers
	}{
		{
			`nil`,
			nil,
			nil,
		},
		{
			`empty slice`,
			resolvers{},
			resolvers{},
		},
		{
			`one element`,
			resolvers{a},
			resolvers{a},
		},
		{
			`two elements`,
			resolvers{a, b},
			resolvers{a, b},
		},
		{
			`three elements`,
			resolvers{a, b, c},
			resolvers{a, b, c},
		},
		{
			`mixed nils`,
			resolvers{a, nil, nil, b, c, nil},
			resolvers{a, b, c},
		},
		{
			`nested`,
			resolvers{resolvers{a, b}, c},
			resolvers{a, b, c},
		},
		{
			`nested with nils`,
			resolvers{resolvers{nil, a, b}, nil, nil, c, nil},
			resolvers{a, b, c},
		},
		{
			`deeply nested`,
			resolvers{resolvers{resolvers{resolvers{a}, b}}, c},
			resolvers{a, b, c},
		},
		{
			`deeply nested with nils`,
			resolvers{nil, resolvers{resolvers{nil, resolvers{a, nil, nil}, b}}, resolvers{}, nil, c, nil},
			resolvers{a, b, c},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.args...)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Factories()\nGot:  %#v\nWant: %#v", got, tt.want)
			}
		})
	}
}
