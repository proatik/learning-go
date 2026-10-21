package clean

import (
	"slices"
	"testing"
)

func TestRemoveNegatives(t *testing.T) {
	tests := []struct {
		nums []int
		want []int
	}{
		{[]int{3, -1, 0, -7, 5}, []int{3, 0, 5}},
		{[]int{-1, -2}, []int{}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{0}, []int{0}},
		{[]int{-5, 4, -3, 2, -1}, []int{4, 2}},
		{[]int{9, -9, 9}, []int{9, 9}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := RemoveNegatives(tt.nums)
		if !slices.Equal(got, tt.want) {
			t.Errorf("RemoveNegatives(%v) = %v, want %v", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("RemoveNegatives modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestReversed(t *testing.T) {
	tests := []struct {
		words []string
		want  []string
	}{
		{[]string{"a", "b", "c"}, []string{"c", "b", "a"}},
		{[]string{"go", "zig"}, []string{"zig", "go"}},
		{[]string{"solo"}, []string{"solo"}},
		{[]string{"x", "y", "y", "z"}, []string{"z", "y", "y", "x"}},
		{nil, []string{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.words)
		got := Reversed(tt.words)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Reversed(%q) = %q, want %q", orig, got, tt.want)
		}
		if !slices.Equal(tt.words, orig) {
			t.Errorf("Reversed modified its input: %q became %q", orig, tt.words)
		}
	}
}

func TestReversedIsACopy(t *testing.T) {
	words := []string{"a", "b"}
	got := Reversed(words)
	if len(got) != 2 {
		t.Fatalf("Reversed(%q) = %q, want 2 elements", words, got)
	}
	got[0] = "changed"
	if words[0] != "a" || words[1] != "b" {
		t.Errorf("changing the result changed the input: %q", words)
	}
}

func TestDescending(t *testing.T) {
	tests := []struct {
		nums []int
		want []int
	}{
		{[]int{2, 9, 4, 9}, []int{9, 9, 4, 2}},
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{3, 2, 1}, []int{3, 2, 1}},
		{[]int{-1, 5, 0, -8}, []int{5, 0, -1, -8}},
		{[]int{7}, []int{7}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := Descending(tt.nums)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Descending(%v) = %v, want %v", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("Descending modified its input: %v became %v", orig, tt.nums)
		}
	}
}
