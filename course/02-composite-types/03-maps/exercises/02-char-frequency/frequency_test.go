package frequency

import (
	"maps"
	"testing"
)

func TestCharFrequency(t *testing.T) {
	tests := []struct {
		in   string
		want map[rune]int
	}{
		{"hello", map[rune]int{'h': 1, 'e': 1, 'l': 2, 'o': 1}},
		{"a a A", map[rune]int{'a': 2, 'A': 1}},
		{"héé", map[rune]int{'h': 1, 'é': 2}},
		{"日本日", map[rune]int{'日': 2, '本': 1}},
		{"x\ty\nx", map[rune]int{'x': 2, 'y': 1}},
		{"!!??", map[rune]int{'!': 2, '?': 2}},
		{"", map[rune]int{}},
		{"  \n", map[rune]int{}},
	}
	for _, tt := range tests {
		got := CharFrequency(tt.in)
		if got == nil {
			t.Errorf("CharFrequency(%q) = nil, want an initialized map", tt.in)
			continue
		}
		if !maps.Equal(got, tt.want) {
			t.Errorf("CharFrequency(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestFirstUnique(t *testing.T) {
	tests := []struct {
		in     string
		want   rune
		wantOK bool
	}{
		{"swiss", 'w', true},
		{"aabbc", 'c', true},
		{"abc", 'a', true},
		{"héllo hé", 'o', true},
		{"aa bb", ' ', true},
		{"日本日", '本', true},
		{"zyzy x", ' ', true},
		{"aabb", 0, false},
		{"", 0, false},
		{"abab  ", 0, false},
	}
	for _, tt := range tests {
		got, ok := FirstUnique(tt.in)
		if got != tt.want || ok != tt.wantOK {
			t.Errorf("FirstUnique(%q) = (%q, %t), want (%q, %t)", tt.in, got, ok, tt.want, tt.wantOK)
		}
	}
}
