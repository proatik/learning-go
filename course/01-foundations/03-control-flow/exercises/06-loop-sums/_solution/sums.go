package sums

// SumTo returns 1 + 2 + … + n, or 0 when n <= 0.
func SumTo(n int) int {
	sum := 0
	for i := range n + 1 { // 0 … n; adding 0 is harmless
		sum += i
	}
	return sum
}

// SumOfMultiples returns the sum of all k with 1 <= k < limit that are
// divisible by a or by b. It assumes a > 0 and b > 0.
func SumOfMultiples(limit, a, b int) int {
	sum := 0
	for k := range limit { // 0 … limit-1
		if k == 0 || (k%a != 0 && k%b != 0) {
			continue // not a multiple we want: skip to the next k
		}
		sum += k
	}
	return sum
}

// CountDivisors returns how many integers d with 1 <= d <= n divide n
// exactly. It returns 0 when n <= 0.
func CountDivisors(n int) int {
	count := 0
	for i := range n {
		if n%(i+1) == 0 { // candidate divisors are 1 … n
			count++
		}
	}
	return count
}
