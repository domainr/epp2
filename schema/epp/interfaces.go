package epp

import "github.com/domainr/epp2/internal/xml"

// Body represents a valid EPP body element.
//
// Standard EPP body elements include <hello>, <greeting>, <command>, and
// <response>.
type Body interface {
	eppBody()
}

// Action is a generic EPP command action.
//
// An Action is serialized to XML as the first child of a <command> element.
type Action interface {
	eppAction()
}

type NamedCommand interface {
	EPPCommandName() xml.Name
}

// CheckType is a child element of EPP <check>.
//
// It is represented as a <check> element with an object-specific namespace.
type CheckType interface {
	EPPCheck()
}

// Value is a generic EPP result value.
//
// It is represented as a <value> element with an object or extension-specific
// namespace.
type Value interface {
	EPPValue()
}

// ResponseData is a generic EPP response data type, serialized as child elements
// of a <resData> (response data) element, containing data specific to the
// command and associated object.
type ResponseData interface {
	EPPResponseData()
}

// Extension is generic EPP extension data.
//
// An <extension> element will contain one or more Extension values.
type Extension interface {
	EPPExtension()
}
