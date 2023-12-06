//go:build !stdxml
// +build !stdxml

package xml

import "github.com/nbio/xml"

type Attr = xml.Attr
type Name = xml.Name
type StartElement = xml.StartElement
type EndElement = xml.EndElement
type Encoder = xml.Encoder
type Decoder = xml.Decoder
type Marshaler = xml.Marshaler
type Unmarshaler = xml.Unmarshaler

var NewEncoder = xml.NewEncoder
var NewDecoder = xml.NewDecoder
var Marshal = xml.Marshal
var Unmarshal = xml.Unmarshal
