package schema

import (
	"encoding"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/domainr/epp2/internal/xml"
)

func Scan(r xml.TokenReader, v interface{}) error {
	v = scanInterface(v)
	textUnmarshaler, _ := v.(encoding.TextUnmarshaler)
	var charData xml.CharData
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

		// Look for other tokens.
		switch v := v.(type) {
		case xml.CharData:
			if len(stack) == 0 {
				if textUnmarshaler != nil {
					charData = append(charData, v...)
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
	ScanStartElement(xml.TokenReader, xml.StartElement) error
}

type EndElementScanner interface {
	ScanEndElement(xml.TokenReader, xml.EndElement) error
}

type EndElementError xml.EndElement

func (e EndElementError) Error() string {
	return "unbalanced end tag: " + e.Name.Local
}

func scanInterface(v interface{}) interface{} {
	switch v := v.(type) {
	case *int:
		return (*Int)(v)
	case *int8:
		return (*Int8)(v)
	case *int16:
		return (*Int16)(v)
	case *int32:
		return (*Int32)(v)
	case *int64:
		return (*Int64)(v)
	case *uint:
		return (*Uint)(v)
	case *uint8:
		return (*Uint8)(v)
	case *uint16:
		return (*Uint16)(v)
	case *uint32:
		return (*Uint32)(v)
	case *uint64:
		return (*Uint64)(v)
	case *float32:
		return (*Float32)(v)
	case *float64:
		return (*Float64)(v)
	case *[]byte:
		return (*ByteSlice)(v)
	case *string:
		return (*String)(v)
	default:
		return v
	}
}

// Int is an int value that implements encoding.TextUnmarshaler.
type Int int

func (v *Int) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = Int(i)
	return nil
}

// Int8 is an int8 value that implements encoding.TextUnmarshaler.
type Int8 int8

func (v *Int8) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 8)
	if err != nil {
		return err
	}
	*v = Int8(i)
	return nil
}

// Int16 is an int16 value that implements encoding.TextUnmarshaler.
type Int16 int16

func (v *Int16) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 16)
	if err != nil {
		return err
	}
	*v = Int16(i)
	return nil
}

// Int32 is an int32 value that implements encoding.TextUnmarshaler.
type Int32 int32

func (v *Int32) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 32)
	if err != nil {
		return err
	}
	*v = Int32(i)
	return nil
}

// Int64 is an int64 value that implements encoding.TextUnmarshaler.
type Int64 int64

func (v *Int64) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		return err
	}
	*v = Int64(i)
	return nil
}

// Uint is a uint value that implements encoding.TextUnmarshaler.
type Uint uint

func (v *Uint) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(strings.TrimSpace(string(text)), 10, strconv.IntSize)
	if err != nil {
		return err
	}
	*v = Uint(i)
	return nil
}

// Uint8 is a uint8 value that implements encoding.TextUnmarshaler.
type Uint8 uint8

func (v *Uint8) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(strings.TrimSpace(string(text)), 10, 8)
	if err != nil {
		return err
	}
	*v = Uint8(i)
	return nil
}

// Uint16 is a uint16 value that implements encoding.TextUnmarshaler.
type Uint16 uint16

func (v *Uint16) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(strings.TrimSpace(string(text)), 10, 16)
	if err != nil {
		return err
	}
	*v = Uint16(i)
	return nil
}

// Uint32 is a uint32 value that implements encoding.TextUnmarshaler.
type Uint32 uint32

func (v *Uint32) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(strings.TrimSpace(string(text)), 10, 32)
	if err != nil {
		return err
	}
	*v = Uint32(i)
	return nil
}

// Uint64 is a uint64 value that implements encoding.TextUnmarshaler.
type Uint64 uint64

func (v *Uint64) UnmarshalText(text []byte) error {
	i, err := strconv.ParseUint(strings.TrimSpace(string(text)), 10, 64)
	if err != nil {
		return err
	}
	*v = Uint64(i)
	return nil
}

// Float32 is a float32 value that implements encoding.TextUnmarshaler.
type Float32 float32

func (v *Float32) UnmarshalText(text []byte) error {
	i, err := strconv.ParseFloat(strings.TrimSpace(string(text)), 32)
	if err != nil {
		return err
	}
	*v = Float32(i)
	return nil
}

// Float64 is a float64 value that implements encoding.TextUnmarshaler.
type Float64 float64

func (v *Float64) UnmarshalText(text []byte) error {
	i, err := strconv.ParseFloat(strings.TrimSpace(string(text)), 64)
	if err != nil {
		return err
	}
	*v = Float64(i)
	return nil
}

// ByteSlice is a slice of byte that implements encoding.TextUnmarshaler.
type ByteSlice []byte

func (s *ByteSlice) UnmarshalText(text []byte) error {
	*s = ByteSlice(text[:])
	return nil
}

// String is a string value that implements encoding.TextUnmarshaler.
type String string

func (s *String) UnmarshalText(text []byte) error {
	*s = String(text)
	return nil
}
