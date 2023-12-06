package epp

// Transfer represents an EPP <transfer> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.4.
type Transfer struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 transfer"`
	// TODO: finish this.
}

func (Transfer) EPPAction() string { return "transfer" }
