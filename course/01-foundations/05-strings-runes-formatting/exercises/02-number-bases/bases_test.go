package bases

import "testing"

func TestToBase(t *testing.T) {
	tests := []struct {
		n, base int
		want    string
	}{
		{10, 2, "1010"},
		{0, 2, "0"},
		{255, 16, "ff"},
		{-8, 8, "-10"},
		{35, 36, "z"},
		{1000, 10, "1000"},
	}
	for _, tt := range tests {
		if got := ToBase(tt.n, tt.base); got != tt.want {
			t.Errorf("ToBase(%d, %d) = %q, want %q", tt.n, tt.base, got, tt.want)
		}
	}
}

func TestFromBase(t *testing.T) {
	tests := []struct {
		s    string
		base int
		want int
	}{
		{"ff", 16, 255},
		{"FF", 16, 255},
		{"1010", 2, 10},
		{"-10", 8, -8},
		{"z", 36, 35},
		{"0", 10, 0},
	}
	for _, tt := range tests {
		got, err := FromBase(tt.s, tt.base)
		if err != nil || got != tt.want {
			t.Errorf("FromBase(%q, %d) = (%d, %v), want (%d, nil)", tt.s, tt.base, got, err, tt.want)
		}
	}
}

func TestFromBaseErrors(t *testing.T) {
	tests := []struct {
		s    string
		base int
	}{
		{"12", 2},
		{"", 10},
		{"g", 16},
		{"1.5", 10},
		{"0x1f", 16},
	}
	for _, tt := range tests {
		got, err := FromBase(tt.s, tt.base)
		if err == nil || got != 0 {
			t.Errorf("FromBase(%q, %d) = (%d, %v), want (0, error)", tt.s, tt.base, got, err)
		}
	}
}

func TestDescribe(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{255, "255 = 0b11111111 = 0o377 = 0xFF"},
		{0, "0 = 0b0 = 0o0 = 0x0"},
		{10, "10 = 0b1010 = 0o12 = 0xA"},
		{4096, "4096 = 0b1000000000000 = 0o10000 = 0x1000"},
	}
	for _, tt := range tests {
		if got := Describe(tt.n); got != tt.want {
			t.Errorf("Describe(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
