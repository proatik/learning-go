package iterate

// CollatzSteps returns how many Collatz steps it takes for n to reach 1:
// halve n when it is even, otherwise replace it with 3n+1.
// It returns -1 when n < 1.
func CollatzSteps(n int) int {
	if n < 1 {
		return -1 // 0 and negatives never reach 1
	}
	steps := 0
	for n != 1 { // Go's "while" loop
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		steps++
	}
	return steps
}

// GCD returns the greatest common divisor of a and b, which is never
// negative. GCD(0, 0) is 0.
func GCD(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b // both right-hand values are computed before assigning
	}
	return a
}
