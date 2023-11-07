package schema

import "github.com/domainr/epp2/internal/xml"

// Any represents an arbitrary XML tag and its contents. It is used when
// unmarshaling with a [Resolver] to represent unrecognized elements.
type Any struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	InnerXML string     `xml:",innerxml"`
}
