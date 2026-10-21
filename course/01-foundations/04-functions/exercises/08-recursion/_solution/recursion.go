package recursion

import "fmt"

// Factorial returns n!. It returns an error if n is negative or
// if n! does not fit in an int (n > 20).
func Factorial(n int) (int, error) {
	if n < 0 {
		return 0, fmt.Errorf("factorial of negative number %d", n)
	}
	if n > 20 {
		return 0, fmt.Errorf("factorial of %d overflows int", n)
	}
	if n == 0 {
		return 1, nil // base case
	}
	f, err := Factorial(n - 1)
	if err != nil {
		return 0, err
	}
	return n * f, nil
}

// Power returns base raised to exp (exp >= 0) using exponentiation by squaring.
func Power(base, exp int) int {
	if exp == 0 {
		return 1
	}
	if exp%2 == 0 {
		half := Power(base, exp/2) // compute once, then square
		return half * half
	}
	return base * Power(base, exp-1)
}

// GCD returns the greatest common divisor of a and b (a, b >= 0).
func GCD(a, b int) int {
	if b == 0 {
		return a
	}
	return GCD(b, a%b)
}

// SumDigits returns the sum of the decimal digits of n.
// The sign of n is ignored.
func SumDigits(n int) int {
	if n < 0 {
		return SumDigits(-n)
	}
	if n < 10 {
		return n
	}
	return n%10 + SumDigits(n/10)
}
