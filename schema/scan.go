package schema

import (
	"encoding"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/domainr/epp2/internal/xml"
)

// ScanFor scans TokenReader r for an XML tag matching name.
// It will then scan the contents of that tag into v.
func ScanFor(r xml.TokenReader, name xml.Name, v interface{}) error {
	return Scan(r, ElementScannerFunc(func(e xml.StartElement) (interface{}, error) {
		if e.Name == name {
			return v, nil
		}
		return nil, nil
	}))
}

// Scan scans TokenReader r into v. It performs a shallow scan, looking for XML
// start tags and character data. If a caller wishes to scan child nodes, then v
// must implement ElementScanner and return a target for Scan to scan child
// nodes into.
//
// If v implements ElementScanner, then Scan will call v.ScanElement and
// recursively call Scan on the returned value.
//
// If v implements CharDataScanner, then Scan will call v.ScanCharData for all
// child CharData tokens.
//
// If v implements encoding.TextUnmarshaler, then Scan will accumulate XML
// character data and call v.UnmarshalText once.
//
// If v is a pointer to a built-in type (e.g. int, string, bool, etc.), Scan will
// attempt to unmarshal the XML character data into v.
func Scan(r xml.TokenReader, v interface{}) error {
	v = scanInterface(v)
	var err error

	var charData xml.CharData
	textUnmarshaler, _ := v.(encoding.TextUnmarshaler)
	if textUnmarshaler != nil {
		defer func() {
			if len(charData) == 0 {
				return
			}
			serr := err
			err = textUnmarshaler.UnmarshalText(charData)
			if err == nil {
				err = serr
			}
		}()
	}

	var name *xml.Name

	for {
		var t xml.Token
		t, err = r.Token()
		if t == nil && err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		// Look for a start element first.
		if start, ok := t.(xml.StartElement); ok {
			name = &start.Name
			var v2 interface{}
			if s, ok := v.(ElementScanner); ok {
				v2, err = s.ScanElement(start)
				if err != nil {
					return err
				}
			}
			err = Scan(r, v2)
			if end, ok := err.(EndElementError); ok {
				t = xml.EndElement(end)
			} else if err != nil {
				return err
			}
		}

		// An unbalanced end element might have been returned from Scan above.
		if end, ok := t.(xml.EndElement); ok {
			if name == nil {
				return EndElementError(end)
			}
			if end.Name != *name {
				return fmt.Errorf("unexpected end tag %s, want %s", end.Name.Local, name.Local)
			}
			name = nil
			continue
		}

		// Look for other tokens.
		switch t := t.(type) {
		case xml.CharData:
			if textUnmarshaler != nil {
				charData = append(charData, t...)
			}
		}
	}
}

type ElementScanner interface {
	ScanElement(xml.StartElement) (interface{}, error)
}

type AttrScanner interface {
	ScanAttr(xml.Attr) (interface{}, error)
}

type CharDataScanner interface {
	ScanCharData(xml.CharData) error
}

type ElementScannerFunc func(xml.StartElement) (interface{}, error)

func (f ElementScannerFunc) ScanElement(e xml.StartElement) (interface{}, error) {
	return f(e)
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
