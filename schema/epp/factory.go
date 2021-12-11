package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// Factory returns a schema.Factory for types defined in the this package.
func Factory() schema.Factory {
	return schema.FactoryFunc(factory)
}

func factory(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	// Body
	case "hello":
		return &Hello{}
	case "greeting":
		return &Greeting{}
	case "command":
		return &Command{}
	case "response":
		return &Response{}

	// CommandType
	case "check":
		return &Check{}
	case "create":
		return &Create{}
	case "delete":
		return &Delete{}
	case "info":
		return &Info{}
	case "login":
		return &Login{}
	case "logout":
		return &Logout{}
	case "poll":
		return &Poll{}
	case "renew":
		return &Renew{}
	case "transfer":
		return &Transfer{}
	case "update":
		return &Update{}
	}
	return nil
}
