package protocol

import (
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/contact"
	"github.com/domainr/epp2/schema/domain"
	"github.com/domainr/epp2/schema/epp"
	"github.com/domainr/epp2/schema/eppcom"
)

// defaultSchemas is an array (not a slice) so DefaultSchemas can return a copy
// that callers can mutate.
var defaultSchemas = [...]schema.Schema{
	epp.Schema,
	eppcom.Schema,
	contact.Schema,
	domain.Schema,
}

// DefaultSchemas returns the default set of [schema.Schema] used by this package.
func DefaultSchemas() schema.Schemas {
	schemas := defaultSchemas
	return schemas[:]
}
