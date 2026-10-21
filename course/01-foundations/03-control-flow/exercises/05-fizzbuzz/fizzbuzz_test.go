package fizzbuzz

import "testing"

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "1"}, {2, "2"}, {3, "Fizz"}, {4, "4"}, {5, "Buzz"},
		{6, "Fizz"}, {9, "Fizz"}, {10, "Buzz"}, {15, "FizzBuzz"},
		{30, "FizzBuzz"}, {45, "FizzBuzz"}, {98, "98"}, {100, "Buzz"},
		{0, "FizzBuzz"}, {-3, "Fizz"}, {-5, "Buzz"}, {-15, "FizzBuzz"}, {-7, "-7"},
	}
	for _, tt := range tests {
		if got := FizzBuzz(tt.n); got != tt.want {
			t.Errorf("FizzBuzz(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func ExamplePrintFizzBuzz() {
	PrintFizzBuzz(15)
	// Output:
	// 1
	// 2
	// Fizz
	// 4
	// Buzz
	// Fizz
	// 7
	// 8
	// Fizz
	// Buzz
	// 11
	// Fizz
	// 13
	// 14
	// FizzBuzz
}

func ExamplePrintFizzBuzz_one() {
	PrintFizzBuzz(1)
	// Output:
	// 1
}

func ExamplePrintFizzBuzz_zero() {
	PrintFizzBuzz(0)
	PrintFizzBuzz(-4)
	// Output:
}
