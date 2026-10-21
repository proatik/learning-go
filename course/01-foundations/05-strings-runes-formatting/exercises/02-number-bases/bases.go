package bases

// ToBase returns n written in the given base (2 to 36), using strconv.FormatInt.
func ToBase(n, base int) string {
	// TODO
	return ""
}

// FromBase parses s as a number in the given base, using strconv.ParseInt.
// It returns an error if s is not a valid number in that base.
func FromBase(s string, base int) (int, error) {
	// TODO
	return 0, nil
}

// Describe returns n in decimal, binary, octal and hexadecimal,
// e.g. "255 = 0b11111111 = 0o377 = 0xFF".
func Describe(n int) string {
	// TODO: one fmt.Sprintf call.
	return ""
}
