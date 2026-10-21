package bylength

import (
	"slices"
	"testing"
)

func TestSortByLength(t *testing.T) {
	tests := []struct {
		words, want []string
	}{
		{[]string{"pear", "fig", "banana", "kiwi", "apple"}, []string{"fig", "kiwi", "pear", "apple", "banana"}},
		{[]string{"bb", "a", "ccc", "aa"}, []string{"a", "aa", "bb", "ccc"}},
		{[]string{"dog", "cat", "ant"}, []string{"ant", "cat", "dog"}},
		{[]string{"Bob", "ann", "Al"}, []string{"Al", "Bob", "ann"}},
		{[]string{"héllo", "hello", "hi"}, []string{"hi", "hello", "héllo"}},
		{[]string{"日本", "abc"}, []string{"日本", "abc"}}, // 2 runes < 3 runes
		{[]string{"x", "x", ""}, []string{"", "x", "x"}},
		{[]string{"only"}, []string{"only"}},
		{nil, nil},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.words)
		SortByLength(tt.words)
		if !slices.Equal(tt.words, tt.want) {
			t.Errorf("SortByLength(%q) gave %q, want %q", orig, tt.words, tt.want)
		}
	}
}

func TestSortByLengthDesc(t *testing.T) {
	tests := []struct {
		words, want []string
	}{
		{[]string{"pear", "fig", "banana", "kiwi", "apple"}, []string{"banana", "apple", "kiwi", "pear", "fig"}},
		{[]string{"bb", "a", "ccc", "aa"}, []string{"ccc", "aa", "bb", "a"}},
		{[]string{"dog", "cat", "ant"}, []string{"ant", "cat", "dog"}},
		{[]string{"a", "bb", "ccc", "dddd"}, []string{"dddd", "ccc", "bb", "a"}},
		{[]string{"日本", "abc"}, []string{"abc", "日本"}},
		{[]string{"z", "zz", "y", "yy"}, []string{"yy", "zz", "y", "z"}},
		{nil, nil},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.words)
		SortByLengthDesc(tt.words)
		if !slices.Equal(tt.words, tt.want) {
			t.Errorf("SortByLengthDesc(%q) gave %q, want %q", orig, tt.words, tt.want)
		}
	}
}

func TestShortest(t *testing.T) {
	tests := []struct {
		words []string
		n     int
		want  []string
	}{
		{[]string{"pear", "fig", "banana", "kiwi"}, 2, []string{"fig", "kiwi"}},
		{[]string{"pear", "fig", "banana", "kiwi"}, 4, []string{"fig", "kiwi", "pear", "banana"}},
		{[]string{"pear", "fig", "banana", "kiwi"}, 1, []string{"fig"}},
		{[]string{"b", "a"}, 5, []string{"a", "b"}},
		{[]string{"b", "a"}, 0, []string{}},
		{[]string{"b", "a"}, -3, []string{}},
		{nil, 2, []string{}},
	}
	for _, tt := range tests {
		orig := slices.Clone(tt.words)
		got := Shortest(tt.words, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Shortest(%q, %d) = %q, want %q", orig, tt.n, got, tt.want)
		}
		if !slices.Equal(tt.words, orig) {
			t.Errorf("Shortest modified its input: %q became %q", orig, tt.words)
		}
	}
}
