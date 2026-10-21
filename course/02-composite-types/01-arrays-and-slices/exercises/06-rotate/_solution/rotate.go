package rotate

// RotateLeft rotates s left by k positions in place: s[k] becomes the first
// element. k may be any int; negative k rotates right.
//
// The three-reversal trick: to turn AB into BA (A = s[:k], B = s[k:]),
// reverse A, reverse B, then reverse the whole slice. Each reversal works on
// a sub-slice that shares s's array, so nothing is allocated.
func RotateLeft(s []int, k int) {
	n := len(s)
	if n == 0 {
		return // also avoids a division by zero in k % n
	}
	k %= n
	if k < 0 {
		k += n // Go's % keeps the sign of the dividend: -1 % 5 == -1
	}
	reverse(s[:k])
	reverse(s[k:])
	reverse(s)
}

// RotateRight rotates s right by k positions in place: the last k elements
// move to the front. k may be any int; negative k rotates left.
func RotateRight(s []int, k int) {
	RotateLeft(s, -k)
}

// reverse reverses s in place.
func reverse(s []int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
