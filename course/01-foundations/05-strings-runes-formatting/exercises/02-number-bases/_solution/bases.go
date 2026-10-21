package bases

import (
	"fmt"
	"strconv"
)

// ToBase returns n written in the given base (2 to 36), using strconv.FormatInt.
func ToBase(n, base int) string {
	return strconv.FormatInt(int64(n), base)
}

// FromBase parses s as a number in the given base, using strconv.ParseInt.
// It returns an error if s is not a valid number in that base.
func FromBase(s string, base int) (int, error) {
	v, err := strconv.ParseInt(s, base, 64)
	if err != nil {
		return 0, err // e.g. strconv.ParseInt: parsing "12": invalid syntax
	}
	return int(v), nil
}

// Describe returns n in decimal, binary, octal and hexadecimal,
// e.g. "255 = 0b11111111 = 0o377 = 0xFF".
func Describe(n int) string {
	// %[1] makes every verb reuse the first argument.
	return fmt.Sprintf("%[1]d = 0b%[1]b = 0o%[1]o = 0x%[1]X", n)
}
