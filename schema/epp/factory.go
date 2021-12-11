package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp2/schema"
)

// factory maps EPP tag names to Go types.
var factory = schema.FactoryFunc(func(name xml.Name) interface{} {
	if name.Space != NS {
		return nil
	}
	switch name.Local {
	// epp.Body
	case "hello":
		return &Hello{}
	case "greeting":
		return &Greeting{}
	case "command":
		return &Command{}
	case "response":
		return &Response{}

	// epp.CommandType
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
})
