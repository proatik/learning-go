package divide

import "errors"

// Divide returns the quotient and remainder of a / b.
// It returns an error if b is zero.
func Divide(a, b int) (q, r int, err error) {
	if b == 0 {
		// q and r still hold their zero values.
		err = errors.New("division by zero")
		return
	}
	q, r = a/b, a%b
	return
}
