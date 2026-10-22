package joiner

import "testing"

// TestJoinInts checks JoinInts with a table of inputs.
func TestJoinInts(t *testing.T) {
	// TODO: a table of {nums []int; sep, want string} cases.
	// Cover: empty, one number, several numbers, negative numbers,
	// and a separator other than ",".
	// Report: "JoinInts(%v, %q) = %q, want %q"
	_ = JoinInts
}

// BenchmarkJoinInts measures JoinInts on 100 numbers.
func BenchmarkJoinInts(b *testing.B) {
	// TODO: build the input once, before the loop (setup is not measured).
	// TODO: call b.ReportAllocs().
	// TODO: call JoinInts inside: for b.Loop() { ... }
}
