package luhn

// Valid reports whether number passes the Luhn checksum. Spaces are ignored;
// any other non-digit makes it invalid, as do fewer than two digits.
func Valid(number string) bool {
	// TODO
	return false
}

// CheckDigit returns the digit that makes payload+digit a valid Luhn number.
func CheckDigit(payload string) (int, error) {
	// TODO
	return 0, nil
}
