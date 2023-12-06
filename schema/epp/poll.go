package epp

// Poll represents an EPP <poll> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.3.
type Poll struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 poll"`
	// TODO: finish this.
}

func (Poll) EPPAction() string { return "poll" }
