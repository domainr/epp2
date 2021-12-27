package schema

import (
	"fmt"
	"io"

	"github.com/domainr/epp2/internal/xml"
)

func Scan(r xml.TokenReader, v interface{}) error {
	type frame struct {
		name   xml.Name
		parent interface{}
	}

	var stack []frame

	for {
		t, terr := r.Token()
		switch t := t.(type) {
		case xml.StartElement:
			stack = append(stack, frame{t.Name, v})
			if s, ok := v.(StartElementScanner); ok {
				v2, err := s.ScanStartElement(t)
				if err != nil {
					return err
				}
				v = v2
			}

		case xml.EndElement:
			if s, ok := v.(EndElementScanner); ok {
				err := s.ScanEndElement(t)
				if err != nil {
					return err
				}
			}
			if len(stack) == 0 {
				return fmt.Errorf("unexpected end tag %s", t.Name.Local)
			}
			frame := stack[len(stack)-1]
			if frame.name != t.Name {
				return fmt.Errorf("unexpected end tag %s, want %s", t.Name.Local, frame.name.Local)
			}
			stack = stack[:len(stack)-1]
			v = frame.parent
		}
		if terr == io.EOF {
			return nil
		} else if terr != nil {
			return terr
		}
	}
}

type StartElementScanner interface {
	ScanStartElement(xml.StartElement) (interface{}, error)
}

type EndElementScanner interface {
	ScanEndElement(xml.EndElement) error
}
