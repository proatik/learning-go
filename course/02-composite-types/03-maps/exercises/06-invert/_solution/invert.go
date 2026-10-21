package invert

import "slices"

// Positions maps each word to the indices where it occurs in words, in
// increasing order. The result is never nil.
func Positions(words []string) map[string][]int {
	pos := map[string][]int{}
	for i, w := range words {
		pos[w] = append(pos[w], i) // pos[w] is nil the first time: append handles it
	}
	return pos
}

// Invert swaps the keys and values of m. Keys that share a value are
// collected into a slice sorted in ascending order. m is not modified and
// the result is never nil.
func Invert(m map[string]int) map[int][]string {
	out := make(map[int][]string, len(m))
	for k, v := range m { // random order...
		out[v] = append(out[v], k)
	}
	for _, keys := range out { // ...so make each group deterministic
		slices.Sort(keys)
	}
	return out
}
