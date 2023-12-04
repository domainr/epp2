package dataunit

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadWrite(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"zero-length data", ""},
		{"one char", "a"},
		{"big value", strings.Repeat("hiya", 1_000_000)},
		{"Unicode", "scholß"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b bytes.Buffer
			err := WriteDataUnit(&b, []byte(tt.data))
			if err != nil {
				t.Errorf("WriteDataUnit: err == %v", err)
			}
			res, err := ReadDataUnit(&b)
			if err != nil {
				t.Errorf("ReadDataUnit: err == %v", err)
			}
			if !bytes.Equal(res, []byte(tt.data)) {
				t.Errorf("ReadDataUnit: got %q, expected %q", string(res), tt.data)
			}
		})
	}
}
