package port

import (
	"errors"
	"testing"
)

// TestParsePort checks valid ports and every error path of ParsePort.
func TestParsePort(t *testing.T) {
	tests := []struct {
		in      string
		want    int
		wantErr error
	}{
		{"80", 80, nil},
		{"1", 1, nil},
		{"65535", 65535, nil},
		{"", 0, ErrEmpty},
		{"http", 0, ErrNotNumber},
		{"8o", 0, ErrNotNumber},
		{"+80", 0, ErrNotNumber},
		{"-1", 0, ErrNotNumber},
		{" 80", 0, ErrNotNumber},
		{"0", 0, ErrOutOfRange},
		{"65536", 0, ErrOutOfRange},
		{"99999999999999999999", 0, ErrOutOfRange},
	}
	for _, tt := range tests {
		got, err := ParsePort(tt.in)
		if !errors.Is(err, tt.wantErr) {
			t.Errorf("ParsePort(%q) error = %v, want %v", tt.in, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParsePort(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}
