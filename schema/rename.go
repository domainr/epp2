package schema

import "github.com/domainr/epp2/internal/xml"

// Rename renames an [xml.StartElement].
func Rename(e xml.StartElement, space, local string) xml.StartElement {
	e.Name.Space = space
	e.Name.Local = local
	return e
}
