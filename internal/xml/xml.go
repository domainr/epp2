//go:build !stdxml
// +build !stdxml

package xml

import "github.com/nbio/xml"

type Token = xml.Token
type Attr = xml.Attr
type Name = xml.Name
type StartElement = xml.StartElement
type EndElement = xml.EndElement
type Encoder = xml.Encoder
type Decoder = xml.Decoder
type TokenReader = xml.TokenReader

var NewEncoder = xml.NewEncoder
var NewDecoder = xml.NewDecoder
var Marshal = xml.Marshal
var Unmarshal = xml.Unmarshal
