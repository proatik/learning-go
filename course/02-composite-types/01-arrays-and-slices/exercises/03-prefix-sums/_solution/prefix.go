package prefix

// PrefixSums returns a new slice p with len(p) == len(s) where
// p[i] = s[0] + s[1] + ... + s[i]. s is not modified.
func PrefixSums(s []int) []int {
	p := make([]int, len(s)) // one allocation of the exact size
	total := 0
	for i, v := range s {
		total += v
		p[i] = total
	}
	return p
}

// RangeSum returns the sum of s[lo:hi], where prefix == PrefixSums(s).
// It assumes 0 <= lo <= hi <= len(prefix).
//
// prefix[hi-1] is the sum of s[0:hi] and prefix[lo-1] is the sum of s[0:lo];
// their difference is the sum of s[lo:hi].
func RangeSum(prefix []int, lo, hi int) int {
	if lo == hi {
		return 0
	}
	sum := prefix[hi-1]
	if lo > 0 {
		sum -= prefix[lo-1]
	}
	return sum
}
