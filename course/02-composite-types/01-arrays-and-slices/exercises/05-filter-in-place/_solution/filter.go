package filter

// KeepEven returns the even elements of s, in order, reusing s's backing
// array. The unused tail of s is zeroed.
//
// out starts as s[:0]: length 0 but the same array and capacity as s, so
// append writes over s from the front and never allocates. The write
// position never overtakes the read position, so nothing unread is lost.
func KeepEven(s []int) []int {
	out := s[:0]
	for _, v := range s {
		if v%2 == 0 {
			out = append(out, v)
		}
	}
	clear(s[len(out):]) // zero the stale leftovers
	return out
}

// RemoveValue returns s without any element equal to v, in order, reusing
// s's backing array. The unused tail of s is zeroed.
func RemoveValue(s []int, v int) []int {
	out := s[:0]
	for _, x := range s {
		if x != v {
			out = append(out, x)
		}
	}
	clear(s[len(out):])
	return out
}
