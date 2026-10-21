package luhn

import "errors"

// Valid reports whether number passes the Luhn checksum. Spaces are ignored;
// any other non-digit makes it invalid, as do fewer than two digits.
func Valid(number string) bool {
	sum, digits, ok := luhnSum(number, false)
	return ok && digits >= 2 && sum%10 == 0
}

// CheckDigit returns the digit that makes payload+digit a valid Luhn number.
func CheckDigit(payload string) (int, error) {
	// The check digit will be the rightmost digit, so every payload digit is
	// shifted one place left: start doubling at the payload's last digit.
	sum, digits, ok := luhnSum(payload, true)
	if !ok {
		return 0, errors.New("luhn: payload contains a non-digit")
	}
	if digits == 0 {
		return 0, errors.New("luhn: empty payload")
	}
	return (10 - sum%10) % 10, nil
}

// luhnSum walks s from right to left and returns the Luhn sum, the number of
// digits seen, and false if s contains a character other than a digit or space.
// If doubleFirst is true, the rightmost digit is doubled.
func luhnSum(s string, doubleFirst bool) (sum, digits int, ok bool) {
	double := doubleFirst
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c == ' ' {
			continue
		}
		if c < '0' || c > '9' {
			return 0, 0, false
		}
		d := int(c - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		digits++
		double = !double
	}
	return sum, digits, true
}
