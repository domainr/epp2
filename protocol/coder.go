package protocol

import (
	"github.com/domainr/epp2/internal/xml"
	"github.com/domainr/epp2/schema"
	"github.com/domainr/epp2/schema/epp"
)

type coder struct {
	schemas schema.Schemas
}

func (c *coder) marshal(body epp.Body) ([]byte, error) {
	e := epp.EPP{Body: body}
	return xml.Marshal(&e)
}

func (c *coder) unmarshal(data []byte) (epp.Body, error) {
	var e epp.EPP
	err := schema.Unmarshal(data, &e, c.schemas)
	return e.Body, err
}
