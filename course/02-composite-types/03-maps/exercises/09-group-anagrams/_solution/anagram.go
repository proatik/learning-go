package anagram

import (
	"cmp"
	"slices"
	"strings"
)

// Key returns the canonical anagram key of word: its lowercased runes in
// sorted order.
func Key(word string) string {
	rs := []rune(strings.ToLower(word))
	slices.Sort(rs)
	return string(rs) // a string is comparable, so it can be a map key
}

// GroupAnagrams groups words whose keys are equal. Each group is sorted
// ascending, and the groups are ordered by their first word. words is not
// modified.
func GroupAnagrams(words []string) [][]string {
	groups := map[string][]string{}
	for _, w := range words {
		k := Key(w)
		groups[k] = append(groups[k], w)
	}

	out := make([][]string, 0, len(groups))
	for _, g := range groups { // random order: sorted below
		slices.Sort(g)
		out = append(out, g)
	}
	slices.SortFunc(out, func(a, b []string) int {
		return cmp.Compare(a[0], b[0]) // groups are never empty
	})
	return out
}
