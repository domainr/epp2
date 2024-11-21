package epp

// Delete represents an EPP <delete> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.3.1.
type Delete struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 delete"`
	// TODO: finish this.
}

func (Delete) EPPAction() string { return "delete" }
