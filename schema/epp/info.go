package epp

// Info represents an EPP <info> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.2.
type Info struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 info"`
	// TODO: InfoType
}

func (Info) eppCommand() {}
