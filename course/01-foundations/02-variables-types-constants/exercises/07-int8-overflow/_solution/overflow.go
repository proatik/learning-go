package overflow

import "math"

// AddInt8 returns a+b with int8 wraparound, and whether the true sum was
// outside the int8 range.
func AddInt8(a, b int8) (sum int8, overflow bool) {
	// Compute the exact sum in a wider type that cannot overflow.
	wide := int16(a) + int16(b)
	// Converting back to int8 keeps the low 8 bits: that is the wrapped sum.
	return int8(wide), wide > math.MaxInt8 || wide < math.MinInt8
}

// SaturatingAdd returns a+b, clamped to the int8 range instead of wrapping.
func SaturatingAdd(a, b int8) int8 {
	wide := int16(a) + int16(b)
	// The untyped constants take the type int16 from wide.
	return int8(max(math.MinInt8, min(math.MaxInt8, wide)))
}
