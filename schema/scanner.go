package schema

import (
	"io"

	"github.com/domainr/epp2/internal/xml"
)

type Scanner interface {
	Scan(interface{}) error
	TokenReader() xml.TokenReader
}

type scanner struct {
	r    xml.TokenReader
	tags []xml.Name
}

func NewScanner(r xml.TokenReader) Scanner {
	return &scanner{
		r: r,
	}
}

func (s *scanner) TokenReader() xml.TokenReader {
	return s.r
}

func (s *scanner) Scan(v interface{}) error {
	for {
		t, terr := s.r.Token()
		if t != nil {
			switch t := t.(type) {
			case xml.StartElement:
				s.tags = append(s.tags, t.Name)
				depth := len(s.tags)
				if ts, ok := v.(StartElementScanner); ok {
					err := ts.ScanStartElement(s, t)
					if err != nil {
						return err
					}
				}
			}
		}
		if terr == io.EOF {
			return nil
		} else if terr != nil {
			return terr
		}
	}
}

type StartElementScanner interface {
	ScanStartElement(Scanner, xml.StartElement) error
}

type EndElementScanner interface {
	ScanEndElement(Scanner, xml.EndElement) error
}

type Login struct {
	Username    string
	Password    string
	NewPassword *string
}

func (l *Login) ScanStartElement(s Scanner, start xml.StartElement) error {
	switch start.Name.Local {
	case "login":
		return s.Scan(l)
	case "clID":
		return s.Scan(&l.Username)
	case "pw":
		return s.Scan(&l.Password)
	case "newPW":
		l.NewPassword = new(string)
		return s.Scan(l.NewPassword)
	}
	return nil
}
