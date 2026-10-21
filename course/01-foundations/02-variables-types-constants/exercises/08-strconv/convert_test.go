package convert

import "testing"

func TestDouble(t *testing.T) {
	tests := []struct {
		s       string
		want    int
		wantErr bool
	}{
		{"21", 42, false},
		{"0", 0, false},
		{"-7", -14, false},
		{"+5", 10, false},
		{"1000000", 2000000, false},
		{"abc", 0, true},
		{"", 0, true},
		{"1.5", 0, true},
		{" 3", 0, true},
		{"12abc", 0, true},
	}
	for _, tt := range tests {
		got, err := Double(tt.s)
		if (err != nil) != tt.wantErr {
			t.Errorf("Double(%q) error = %v, want error: %v", tt.s, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("Double(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}

func TestAddStrings(t *testing.T) {
	tests := []struct {
		a, b    string
		want    float64
		wantErr bool
	}{
		{"1.5", "2.25", 3.75, false},
		{"1", "2", 3, false},
		{"-0.5", "0.5", 0, false},
		{"1e3", "1", 1001, false},
		{"1", "x", 0, true},
		{"x", "1", 0, true},
		{"", "", 0, true},
		{"7", "", 0, true},
	}
	for _, tt := range tests {
		got, err := AddStrings(tt.a, tt.b)
		if (err != nil) != tt.wantErr {
			t.Errorf("AddStrings(%q, %q) error = %v, want error: %v", tt.a, tt.b, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("AddStrings(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestLabel(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{{7, "item-7"}, {0, "item-0"}, {42, "item-42"}, {-3, "item--3"}, {65, "item-65"}}
	for _, tt := range tests {
		if got := Label(tt.n); got != tt.want {
			t.Errorf("Label(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestBinary(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{{10, "1010"}, {0, "0"}, {1, "1"}, {255, "11111111"}, {256, "100000000"}, {-5, "-101"}}
	for _, tt := range tests {
		if got := Binary(tt.n); got != tt.want {
			t.Errorf("Binary(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
