package overflow

// AddInt8 returns a+b with int8 wraparound, and whether the true sum was
// outside the int8 range.
func AddInt8(a, b int8) (sum int8, overflow bool) {
	return 0, false // TODO
}

// SaturatingAdd returns a+b, clamped to the int8 range instead of wrapping.
func SaturatingAdd(a, b int8) int8 {
	return 0 // TODO
}
