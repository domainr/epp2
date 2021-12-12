package schema

import "github.com/nbio/xml"

// Any represents an arbitrary XML tag and its contents. It is used when
// unmarshaling with a Factory to represent unrecognized elements.
type Any struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	InnerXML string     `xml:",innerxml"`
}
