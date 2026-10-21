package search

// FirstPair returns the first pair (a, b) with 1 <= a <= b <= limit and
// a*b == target, searching a upwards and then b upwards from a.
// It returns (0, 0) if no such pair exists.
func FirstPair(target, limit int) (int, int) {
	x, y := 0, 0 // the answer; stays (0, 0) if nothing is found
outer:
	for a := 1; a <= limit; a++ {
		for b := a; b <= limit; b++ {
			if a*b == target {
				x, y = a, b
				break outer // a plain break would only leave the inner loop
			}
		}
	}
	return x, y
}

// CountPrimes returns how many prime numbers p satisfy 2 <= p <= n.
func CountPrimes(n int) int {
	count := 0
candidates:
	for p := 2; p <= n; p++ {
		for d := 2; d*d <= p; d++ {
			if p%d == 0 {
				continue candidates // p has a divisor: not prime, try the next p
			}
		}
		count++ // the inner loop found no divisor
	}
	return count
}
