//go:build !stdxml
// +build !stdxml

package xml

import "github.com/nbio/xml"

type Name = xml.Name
type Attr = xml.Attr
type Token = xml.Token
type StartElement = xml.StartElement
type EndElement = xml.EndElement
type CharData = xml.CharData
type Encoder = xml.Encoder
type Decoder = xml.Decoder
type TokenReader = xml.TokenReader

var NewEncoder = xml.NewEncoder
var NewDecoder = xml.NewDecoder
var Marshal = xml.Marshal
var Unmarshal = xml.Unmarshal
