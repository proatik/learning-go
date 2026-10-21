package config

import "testing"

func TestParsePort(t *testing.T) {
	tests := []struct {
		s       string
		want    int
		wantErr bool
	}{
		{"", 8080, false},
		{"443", 443, false},
		{"80", 80, false},
		{"1", 1, false},
		{"65535", 65535, false},
		{"8080", 8080, false},
		{"0", 0, true},
		{"65536", 0, true},
		{"-1", 0, true},
		{"http", 0, true},
		{"80.5", 0, true},
	}
	for _, tt := range tests {
		got, err := ParsePort(tt.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParsePort(%q) error = %v, want error: %v", tt.s, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("ParsePort(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestScale(t *testing.T) {
	tests := []struct {
		value   int
		factor  string
		want    int
		wantErr bool
	}{
		{5, "", 5, false},
		{5, "1", 5, false},
		{5, "3", 15, false},
		{5, "-2", -10, false},
		{0, "7", 0, false},
		{-4, "", -4, false},
		{5, "0", 0, true},
		{5, "x", 0, true},
	}
	for _, tt := range tests {
		got, err := Scale(tt.value, tt.factor)
		if (err != nil) != tt.wantErr {
			t.Errorf("Scale(%d, %q) error = %v, want error: %v", tt.value, tt.factor, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Scale(%d, %q) = %d, want %d", tt.value, tt.factor, got, tt.want)
		}
	}
}
