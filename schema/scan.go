package schema

import (
	"fmt"
	"io"

	"github.com/domainr/epp2/internal/xml"
)

func Scan(r xml.TokenReader, v interface{}) error {
	var stack []xml.Name

	for {
		t, terr := r.Token()

		// Look for a start element first.
		if start, ok := t.(xml.StartElement); ok {
			stack = append(stack, start.Name)
			if len(stack) == 1 {
				if s, ok := v.(StartElementScanner); ok {
					err := s.ScanStartElement(r, start)
					if end, ok := err.(EndElementError); ok {
						t = xml.EndElement(end)
					} else if err != nil {
						return err
					}
				}
			}
		}

		// An unbalanced end element might have been returned from ScanStartElement above.
		if end, ok := t.(xml.EndElement); ok {
			if len(stack) == 0 {
				return EndElementError(end)
			}
			name := stack[len(stack)-1]
			if name != end.Name {
				return fmt.Errorf("unexpected end tag %s, want %s", end.Name.Local, name.Local)
			}
			if len(stack) == 1 {
				if s, ok := v.(EndElementScanner); ok {
					err := s.ScanEndElement(r, end)
					if err != nil {
						return err
					}
				}
			}
			stack = stack[:len(stack)-1]
		}

		if terr == io.EOF {
			return nil
		} else if terr != nil {
			return terr
		}
	}
}

type StartElementScanner interface {
	ScanStartElement(xml.TokenReader, xml.StartElement) error
}

type EndElementScanner interface {
	ScanEndElement(xml.TokenReader, xml.EndElement) error
}

type EndElementError xml.EndElement

func (e EndElementError) Error() string {
	return "unbalanced end tag: " + e.Name.Local
}
