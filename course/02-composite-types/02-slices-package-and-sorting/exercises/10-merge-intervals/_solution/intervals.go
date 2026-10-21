package intervals

import (
	"cmp"
	"slices"
)

// Merge returns the union of intervals as non-overlapping intervals sorted
// by start. Overlapping or touching intervals are merged. Each interval is
// {start, end} with start <= end. intervals is not modified.
func Merge(intervals [][2]int) [][2]int {
	// [2]int is an array value, so Clone copies every interval: a full copy.
	sorted := slices.Clone(intervals)
	slices.SortFunc(sorted, func(a, b [2]int) int {
		return cmp.Or(cmp.Compare(a[0], b[0]), cmp.Compare(a[1], b[1]))
	})

	var merged [][2]int
	for _, iv := range sorted {
		last := len(merged) - 1
		if last >= 0 && iv[0] <= merged[last][1] {
			// Overlaps or touches the previous interval: extend it.
			merged[last][1] = max(merged[last][1], iv[1])
			continue
		}
		merged = append(merged, iv)
	}
	return merged
}

// Covered returns the total length of the union of intervals.
func Covered(intervals [][2]int) int {
	total := 0
	for _, iv := range Merge(intervals) {
		total += iv[1] - iv[0]
	}
	return total
}

// Gaps returns the holes between the merged intervals, in order.
func Gaps(intervals [][2]int) [][2]int {
	merged := Merge(intervals)
	var gaps [][2]int
	for i := 1; i < len(merged); i++ {
		// Merged intervals never touch, so each neighbouring pair leaves a hole.
		gaps = append(gaps, [2]int{merged[i-1][1], merged[i][0]})
	}
	return gaps
}
