package epp

// Renew represents an EPP <renew> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.3.1.
type Renew struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 renew"`
	// TODO: finish this.
}

func (Renew) eppCommand() {}
