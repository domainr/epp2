package epp

// Logout represents an EPP <logout> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.1.2.
type Logout struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 logout,selfclosing"`
}

func (Logout) EPPAction() string { return "logout" }
