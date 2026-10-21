package lengths

import "testing"

func TestLengths(t *testing.T) {
	tests := []struct {
		s            string
		bytes, runes int
	}{
		{"", 0, 0},
		{"Go", 2, 2},
		{"héllo", 6, 5},
		{"世界", 6, 2},
		{"🙂", 4, 1},
		{"Hello, 世界!", 14, 10},
		{"a b\tc\n", 6, 6},
	}
	for _, tt := range tests {
		b, r := Lengths(tt.s)
		if b != tt.bytes || r != tt.runes {
			t.Errorf("Lengths(%q) = (%d, %d), want (%d, %d)", tt.s, b, r, tt.bytes, tt.runes)
		}
	}
}

func TestCountVowels(t *testing.T) {
	tests := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"Gopher", 2},
		{"AEIOU aeiou", 10},
		{"rhythm", 0},
		{"héllo, 世界", 1},
		{"Éa", 1},
		{"Programming in Go!", 5},
		{"🙂u🙂", 1},
	}
	for _, tt := range tests {
		if got := CountVowels(tt.s); got != tt.want {
			t.Errorf("CountVowels(%q) = %d, want %d", tt.s, got, tt.want)
		}
	}
}
