package joiner

import "testing"

// TestJoinInts checks JoinInts with a table of inputs.
func TestJoinInts(t *testing.T) {
	tests := []struct {
		nums      []int
		sep, want string
	}{
		{nil, ",", ""},
		{[]int{7}, ",", "7"},
		{[]int{1, 2, 3}, ",", "1,2,3"},
		{[]int{1, 2, 3}, ", ", "1, 2, 3"},
		{[]int{-1, 0, 42}, "|", "-1|0|42"},
		{[]int{10, 20}, "", "1020"},
	}
	for _, tt := range tests {
		if got := JoinInts(tt.nums, tt.sep); got != tt.want {
			t.Errorf("JoinInts(%v, %q) = %q, want %q", tt.nums, tt.sep, got, tt.want)
		}
	}
}

// BenchmarkJoinInts measures JoinInts on 100 numbers.
func BenchmarkJoinInts(b *testing.B) {
	nums := make([]int, 100)
	for i := range nums {
		nums[i] = i * 1000
	}
	b.ReportAllocs()
	for b.Loop() {
		JoinInts(nums, ",")
	}
}
