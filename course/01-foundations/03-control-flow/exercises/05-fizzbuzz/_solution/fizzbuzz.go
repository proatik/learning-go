package fizzbuzz

import (
	"fmt"
	"strconv"
)

// FizzBuzz returns "Fizz" for multiples of 3, "Buzz" for multiples of 5,
// "FizzBuzz" for multiples of both, and the decimal number otherwise.
func FizzBuzz(n int) string {
	switch {
	case n%15 == 0: // divisible by 3 and 5: must come first
		return "FizzBuzz"
	case n%3 == 0:
		return "Fizz"
	case n%5 == 0:
		return "Buzz"
	default:
		return strconv.Itoa(n)
	}
}

// PrintFizzBuzz prints FizzBuzz(i) for i = 1 … n, one per line.
func PrintFizzBuzz(n int) {
	// range n yields 0 … n-1, so shift by one; it runs zero times when n < 1.
	for i := range n {
		fmt.Println(FizzBuzz(i + 1))
	}
}
