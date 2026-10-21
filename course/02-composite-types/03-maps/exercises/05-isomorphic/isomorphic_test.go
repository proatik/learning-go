package isomorphic

import "testing"

func TestIsIsomorphic(t *testing.T) {
	tests := []struct {
		s, t string
		want bool
	}{
		{"egg", "add", true},
		{"paper", "title", true},
		{"foo", "bar", false},
		{"ab", "aa", false},
		{"aa", "ab", false},
		{"badc", "baba", false},
		{"abcabc", "xyzxyz", true},
		{"abcabc", "xyzxzy", false},
		{"日本日", "xyx", true},
		{"xyx", "日本日", true},
		{"héé", "abb", true},
		{"abc", "ab", false},
		{"a", "a", true},
		{"a", "z", true},
		{"", "", true},
		{"", "a", false},
	}
	for _, tt := range tests {
		if got := IsIsomorphic(tt.s, tt.t); got != tt.want {
			t.Errorf("IsIsomorphic(%q, %q) = %t, want %t", tt.s, tt.t, got, tt.want)
		}
	}
}

func TestWordPattern(t *testing.T) {
	tests := []struct {
		pattern, s string
		want       bool
	}{
		{"abba", "dog cat cat dog", true},
		{"abba", "dog cat cat fish", false},
		{"aaaa", "dog cat cat dog", false},
		{"abba", "dog dog dog dog", false},
		{"aaaa", "dog dog dog dog", true},
		{"abc", "dog cat", false},
		{"ab", "dog cat fish", false},
		{"abc", "  red   green blue ", true},
		{"xyx", "go rust go", true},
		{"日本", "sun moon", true},
		{"", "", true},
	}
	for _, tt := range tests {
		if got := WordPattern(tt.pattern, tt.s); got != tt.want {
			t.Errorf("WordPattern(%q, %q) = %t, want %t", tt.pattern, tt.s, got, tt.want)
		}
	}
}
