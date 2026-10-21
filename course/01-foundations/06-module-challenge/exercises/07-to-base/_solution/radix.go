package radix

import "fmt"

// digits[d] is the character for digit value d.
const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

// ToBase formats n in the given base (2 to 36) using digits 0-9a-z.
func ToBase(n int, base int) (string, error) {
	if base < 2 || base > len(digits) {
		return "", fmt.Errorf("radix: base %d is out of range 2..36", base)
	}
	if n == 0 {
		return "0", nil
	}

	// Work with the magnitude as a uint. For math.MinInt, -n overflows
	// back to math.MinInt, but converting first and negating the uint
	// (unsigned arithmetic wraps around) gives the correct magnitude.
	u := uint(n)
	if n < 0 {
		u = -u
	}

	// Repeated division yields the digits least significant first, so
	// each new digit is prepended. Results are at most 65 bytes long,
	// so the simple concatenation is fine.
	b := uint(base)
	s := ""
	for u > 0 {
		s = string(digits[u%b]) + s
		u /= b
	}
	if n < 0 {
		s = "-" + s
	}
	return s, nil
}
