package ranking

import (
	"cmp"
	"fmt"
	"slices"
	"strings"
)

// order returns the indices 0..n-1 sorted by score descending, then name
// ascending. Sorting indices leaves both parallel slices untouched.
func order(names []string, scores []int) ([]int, error) {
	if len(names) != len(scores) {
		return nil, fmt.Errorf("%d names but %d scores", len(names), len(scores))
	}
	idx := make([]int, len(names))
	for i := range idx {
		idx[i] = i
	}
	slices.SortFunc(idx, func(a, b int) int {
		return cmp.Or(
			cmp.Compare(scores[b], scores[a]),   // b, a: higher score first
			strings.Compare(names[a], names[b]), // a, b: alphabetical
		)
	})
	return idx, nil
}

// Rank returns names ordered by score descending, then name ascending.
// names[i] scored scores[i]. It returns an error if the lengths differ.
// Neither input is modified.
func Rank(names []string, scores []int) ([]string, error) {
	idx, err := order(names, scores)
	if err != nil {
		return nil, err
	}
	ranked := make([]string, 0, len(idx))
	for _, i := range idx {
		ranked = append(ranked, names[i])
	}
	return ranked, nil
}

// Standings returns one line "<rank>. <name> (<score>)" per student, in Rank
// order, using competition ranking (1, 1, 3, ...) for equal scores.
// It returns an error if the lengths differ. Neither input is modified.
func Standings(names []string, scores []int) ([]string, error) {
	idx, err := order(names, scores)
	if err != nil {
		return nil, err
	}
	lines := make([]string, 0, len(idx))
	rank := 0
	for pos, i := range idx {
		// A new rank starts whenever the score differs from the previous one;
		// it equals the 1-based position, which skips past earlier ties.
		if pos == 0 || scores[i] != scores[idx[pos-1]] {
			rank = pos + 1
		}
		lines = append(lines, fmt.Sprintf("%d. %s (%d)", rank, names[i], scores[i]))
	}
	return lines, nil
}
