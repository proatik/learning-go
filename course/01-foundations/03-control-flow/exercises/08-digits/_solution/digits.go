package digits

// DigitSum returns the sum of the decimal digits of n, ignoring its sign.
func DigitSum(n int) int {
	if n < 0 {
		n = -n
	}
	sum := 0
	for n > 0 {
		sum += n % 10 // last digit
		n /= 10       // drop it
	}
	return sum
}

// Reverse returns n with its decimal digits in reverse order, keeping the
// sign: Reverse(-120) is -21.
func Reverse(n int) int {
	sign := 1
	if n < 0 {
		sign, n = -1, -n
	}
	rev := 0
	for n > 0 {
		rev = rev*10 + n%10 // shift rev left one digit, append n's last digit
		n /= 10
	}
	return sign * rev
}

// IsPalindrome reports whether n reads the same forwards and backwards.
// Negative numbers are never palindromes.
func IsPalindrome(n int) bool {
	return n >= 0 && Reverse(n) == n
}
