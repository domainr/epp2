package test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/nbio/xml"
)

// RoundTrip validates if v marshals to want or wantErr (if set),
// and the resulting XML unmarshals to v.
func RoundTrip(t *testing.T, v interface{}, wantXML string, wantErr bool) {
	gotXML, err := xml.Marshal(v)
	if (err != nil) != wantErr {
		t.Errorf("xml.Marshal() error = %v, wantErr %v", err, wantErr)
		return
	}
	if string(gotXML) != wantXML {
		t.Errorf("xml.Marshal()\nGot:  %v\nWant: %v", string(gotXML), wantXML)
	}

	if v == nil {
		return
	}

	got := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = xml.Unmarshal(gotXML, got)
	if err != nil {
		t.Errorf("xml.Unmarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(v, got) {
		t.Errorf("xml.Unmarshal()\nGot:  %#v\nWant: %#v", got, v)
	}
}

// RoundTripName validates if v marshals to want or wantErr (if set),
// and the resulting XML unmarshals to v. The outer XML tag will use name, if set.
func RoundTripName(t *testing.T, name xml.Name, v interface{}, want string, wantErr bool) {
	var err error
	buf := &bytes.Buffer{}
	enc := xml.NewEncoder(buf)
	if name == (xml.Name{}) {
		err = enc.Encode(v)
	} else {
		err = enc.EncodeElement(v, xml.StartElement{Name: name})
	}
	if (err != nil) != wantErr {
		t.Errorf("XML encoding error = %v, wantErr %v", err, wantErr)
		return
	}
	if buf.String() != want {
		t.Errorf("XML encoding\nGot:  %v\nWant: %v", buf.String(), want)
	}

	if v == nil {
		return
	}

	got := reflect.New(reflect.TypeOf(v).Elem()).Interface()
	err = xml.Unmarshal(buf.Bytes(), got)
	if err != nil {
		t.Errorf("xml.Unmarshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(v, got) {
		t.Errorf("xml.Unmarshal()\nGot:  %#v\nWant: %#v", got, v)
	}
}
