package epp

// Response represents an EPP server <response> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.6.
type Response struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 response"`

	// Results is one or more <result> elements describing the success or
	// failure of an EPP command.
	Results []Result `xml:"result,omitempty"`

	// MessageQueue is the OPTIONAL <msgQ> element describes messages queued
	// for client retrieval.
	MessageQueue *MessageQueue `xml:"msgQ"`

	// Data is the OPTIONAL <resData> (response data) element
	// contains child elements specific to the command and associated
	// object.
	Data []ResponseData

	// Extensions represents an OPTIONAL <extension> element that MAY
	// be used for server-defined response extensions.
	Extensions Extensions `xml:"extension,omitempty"`

	// TransactionID is a <trID> (transaction identifier) element that
	// contains a CLIENT-generated transaction ID of the command and a
	// SERVER-generated transaction ID that uniquely identifies the
	// response.
	TransactionID TransactionID `xml:"trID"`
}

func (Response) eppBody() {}

// Result represents an EPP server <result> as defined in RFC 5730.
type Result struct {
	Code            ResultCode `xml:"code,attr"`
	Message         Message    `xml:"msg"`
	Values          []Value
	ExtensionValues []ExtensionValue `xml:"extValue,omitempty"`
}

// ExtensionValue wraps an EPP result extension value within a <result>.
type ExtensionValue struct {
	Value  Value
	Reason Message `xml:"reason"`
}

// TransactionID represents an EPP server <trID> as defined in RFC 5730.
type TransactionID struct {
	Client string `xml:"clTRID"`
	Server string `xml:"svTRID"`
}
