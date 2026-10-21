package canon

import (
	"slices"
	"testing"
)

func TestUnique(t *testing.T) {
	tests := []struct {
		nums []int
		want []int
	}{
		{[]int{3, 1, 3, 2, 1}, []int{1, 2, 3}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{5, 5, 5}, []int{5}},
		{[]int{2, 1, 2, 1, 2, 1}, []int{1, 2}},
		{[]int{0, -1, 0, -1}, []int{-1, 0}},
		{[]int{42}, []int{42}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := Unique(tt.nums)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Unique(%v) = %v, want %v", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("Unique modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestCountDistinct(t *testing.T) {
	tests := []struct {
		words []string
		want  int
	}{
		{[]string{"go", "zig", "go"}, 2},
		{[]string{"a", "b", "c"}, 3},
		{[]string{"x", "x", "x", "x"}, 1},
		{[]string{"Go", "go", "GO"}, 3},
		{[]string{"b", "a", "b", "a", ""}, 3},
		{nil, 0},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.words)
		if got := CountDistinct(tt.words); got != tt.want {
			t.Errorf("CountDistinct(%q) = %d, want %d", orig, got, tt.want)
		}
		if !slices.Equal(tt.words, orig) {
			t.Errorf("CountDistinct modified its input: %q became %q", orig, tt.words)
		}
	}
}

func TestSortedLetters(t *testing.T) {
	tests := []struct {
		s, want string
	}{
		{"banana", "aaabnn"},
		{"cba", "abc"},
		{"héllo", "hlloé"},
		{"Go", "Go"},
		{"a b", " ab"},
		{"日本語", "日本語"},
		{"語本日", "日本語"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := SortedLetters(tt.s); got != tt.want {
			t.Errorf("SortedLetters(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}

func TestIsAnagram(t *testing.T) {
	tests := []struct {
		a, b string
		want bool
	}{
		{"listen", "silent", true},
		{"evil", "vile", true},
		{"aab", "abb", false},
		{"abc", "abcd", false},
		{"Tea", "eat", false},
		{"dormitory", "dirtyroom", true},
		{"dormitory", "dirty room", false},
		{"résumé", "mésuré", true},
		{"ab", "ba", true},
		{"", "", true},
		{"a", "", false},
	}
	for _, tt := range tests {
		if got := IsAnagram(tt.a, tt.b); got != tt.want {
			t.Errorf("IsAnagram(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}
