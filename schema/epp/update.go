package epp

// Update represents an EPP <update> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.3.1.
type Update struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 update"`
	// TODO: finish this.
}

func (Update) eppAction() {}
