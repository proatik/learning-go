package seqs

import (
	"slices"
	"testing"
)

func TestEvens(t *testing.T) {
	tests := []struct {
		nums, want []int
	}{
		{[]int{1, 2, 3, 4, -6}, []int{2, 4, -6}},
		{[]int{1, 3, 5}, []int{}},
		{[]int{0, 0}, []int{0, 0}},
		{[]int{8, 6, 7, 5, 4}, []int{8, 6, 4}},
		{[]int{-3, -2, -1}, []int{-2}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := Evens(tt.nums)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Evens(%v) = %v, want %v", orig, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("Evens modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestLastIndex(t *testing.T) {
	tests := []struct {
		words []string
		w     string
		want  int
	}{
		{[]string{"a", "b", "a"}, "a", 2},
		{[]string{"a", "b", "a"}, "b", 1},
		{[]string{"a"}, "z", -1},
		{[]string{"x", "y", "z"}, "x", 0},
		{[]string{"go", "go", "go"}, "go", 2},
		{[]string{"Go"}, "go", -1},
		{nil, "a", -1},
	}
	for _, tt := range tests {
		if got := LastIndex(tt.words, tt.w); got != tt.want {
			t.Errorf("LastIndex(%q, %q) = %d, want %d", tt.words, tt.w, got, tt.want)
		}
	}
}

func TestChunkSums(t *testing.T) {
	tests := []struct {
		nums []int
		n    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 7, 5}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{6, 15}},
		{[]int{1, 2, 3}, 5, []int{6}},
		{[]int{1, 2, 3}, 1, []int{1, 2, 3}},
		{[]int{-1, 1, -2, 2}, 2, []int{0, 0}},
		{[]int{1, 2, 3}, 0, []int{}},
		{[]int{1, 2, 3}, -2, []int{}},
		{nil, 3, []int{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.nums)
		got := ChunkSums(tt.nums, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("ChunkSums(%v, %d) = %v, want %v", orig, tt.n, got, tt.want)
		}
		if !slices.Equal(tt.nums, orig) {
			t.Errorf("ChunkSums modified its input: %v became %v", orig, tt.nums)
		}
	}
}

func TestSortedFields(t *testing.T) {
	tests := []struct {
		s    string
		want []string
	}{
		{"  the quick\tbrown  fox ", []string{"brown", "fox", "quick", "the"}},
		{"c b a", []string{"a", "b", "c"}},
		{"one", []string{"one"}},
		{"b B a A", []string{"A", "B", "a", "b"}},
		{"go\ngo\ngo", []string{"go", "go", "go"}},
		{"   ", []string{}},
		{"", []string{}},
	}
	for _, tt := range tests {
		if got := SortedFields(tt.s); !slices.Equal(got, tt.want) {
			t.Errorf("SortedFields(%q) = %q, want %q", tt.s, got, tt.want)
		}
	}
}
