package dedupe

import (
	"slices"
	"testing"
)

func TestHasDuplicate(t *testing.T) {
	tests := []struct {
		in   []int
		want bool
	}{
		{[]int{1, 2, 3, 1}, true},
		{[]int{1, 2, 3}, false},
		{[]int{7, 7}, true},
		{[]int{0, -1, 1, -1}, true},
		{[]int{5}, false},
		{[]int{0, 0}, true},
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 9}, true},
		{[]int{}, false},
		{nil, false},
	}
	for _, tt := range tests {
		if got := HasDuplicate(tt.in); got != tt.want {
			t.Errorf("HasDuplicate(%v) = %t, want %t", tt.in, got, tt.want)
		}
	}
}

func TestDedupe(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		{[]string{"go", "is", "go", "fun", "is"}, []string{"go", "is", "fun"}},
		{[]string{"a", "A", "a"}, []string{"a", "A"}},
		{[]string{"x", "x", "x"}, []string{"x"}},
		{[]string{"c", "b", "a"}, []string{"c", "b", "a"}},
		{[]string{"", "", "z"}, []string{"", "z"}},
		{[]string{"b", "a", "b", "c", "a", "d"}, []string{"b", "a", "c", "d"}},
		{[]string{}, []string{}},
		{nil, nil},
	}
	for _, tt := range tests {
		in := slices.Clone(tt.in)
		got := Dedupe(in)
		if !slices.Equal(got, tt.want) { // slices.Equal treats nil and empty as equal
			t.Errorf("Dedupe(%q) = %q, want %q", tt.in, got, tt.want)
		}
		if !slices.Equal(in, tt.in) {
			t.Errorf("Dedupe modified its input: %q became %q", tt.in, in)
		}
	}
}
