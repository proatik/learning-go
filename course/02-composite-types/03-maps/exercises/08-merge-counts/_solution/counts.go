package counts

import (
	"fmt"
	"maps"
	"slices"
)

// MergeCounts returns a new map holding, for every key of a or b, the sum
// of its counts. Neither a nor b is modified. The result is never nil.
func MergeCounts(a, b map[string]int) map[string]int {
	// maps.Clone(a) would return nil for a nil a, so allocate explicitly.
	out := make(map[string]int, len(a)+len(b))
	maps.Copy(out, a)
	for k, v := range b {
		out[k] += v // missing key reads as 0
	}
	return out
}

// Prune deletes, in place, every entry of m whose count is less than minCount.
func Prune(m map[string]int, minCount int) {
	maps.DeleteFunc(m, func(_ string, v int) bool { return v < minCount })
}

// PrintCounts prints "key: value" lines sorted by key, or "(none)" if m is
// empty.
func PrintCounts(m map[string]int) {
	if len(m) == 0 { // len of a nil map is 0
		fmt.Println("(none)")
		return
	}
	for _, k := range slices.Sorted(maps.Keys(m)) {
		fmt.Printf("%s: %d\n", k, m[k])
	}
}
