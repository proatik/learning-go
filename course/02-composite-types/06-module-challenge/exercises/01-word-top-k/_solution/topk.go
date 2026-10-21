package topk

import (
	"cmp"
	"slices"
	"strings"
	"unicode"
)

// WordCount is a word together with the number of times it occurs.
type WordCount struct {
	Word  string
	Count int
}

// Counts returns how often each lowercase word occurs in text.
func Counts(text string) map[string]int {
	notWord := func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}
	counts := make(map[string]int)
	for _, w := range strings.FieldsFunc(text, notWord) {
		counts[strings.ToLower(w)]++ // a missing key starts at 0
	}
	return counts
}

// TopK returns the k most frequent words of text, ordered by count
// descending and then word ascending.
func TopK(text string, k int) []WordCount {
	if k <= 0 {
		return []WordCount{}
	}
	counts := Counts(text)

	// Map order is random, so copy into a slice and sort it.
	ws := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		ws = append(ws, WordCount{Word: w, Count: c})
	}
	slices.SortFunc(ws, func(a, b WordCount) int {
		return cmp.Or(
			cmp.Compare(b.Count, a.Count), // descending count
			cmp.Compare(a.Word, b.Word),   // ascending word breaks ties
		)
	})
	return ws[:min(k, len(ws))]
}
