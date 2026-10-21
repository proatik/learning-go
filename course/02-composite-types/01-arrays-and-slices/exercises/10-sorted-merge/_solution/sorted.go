package sorted

// MergeSorted merges the ascending slices a and b into a new ascending slice
// of length len(a)+len(b). Duplicates are kept; a and b are not modified.
//
// i and j point at the next unused element of a and b; each step takes the
// smaller one. When one side runs out, the rest of the other is appended.
func MergeSorted(a, b []int) []int {
	out := make([]int, 0, len(a)+len(b)) // one allocation, never grows
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] <= b[j] { // <= keeps the merge stable: a's equal elements first
			out = append(out, a[i])
			i++
		} else {
			out = append(out, b[j])
			j++
		}
	}
	out = append(out, a[i:]...) // at most one of these two appends anything
	out = append(out, b[j:]...)
	return out
}

// Dedupe removes repeated values from the ascending slice s in place and
// returns the shortened slice, which shares s's backing array.
//
// In a sorted slice equal values are adjacent, so a value is new exactly when
// it differs from the last value kept.
func Dedupe(s []int) []int {
	out := s[:0]
	for _, v := range s {
		if len(out) == 0 || v != out[len(out)-1] {
			out = append(out, v)
		}
	}
	return out
}

// Intersect returns a new ascending slice of the values that appear in both
// ascending slices a and b, each value once.
func Intersect(a, b []int) []int {
	var out []int // a fresh slice: appends never touch a or b
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default: // equal: record it once, then skip every copy on both sides
			v := a[i]
			out = append(out, v)
			for i < len(a) && a[i] == v {
				i++
			}
			for j < len(b) && b[j] == v {
				j++
			}
		}
	}
	return out
}
