package containers

import "testing"

// TestStack checks Stack. Use a table of operations and t.Run subtests.
func TestStack(t *testing.T) {
	// TODO: push/pop order, empty stack, Peek, Len, All order and break.
	_ = NewStack[int]
}

// TestQueue checks Queue. Make it wrap around and grow.
func TestQueue(t *testing.T) {
	// TODO: FIFO order, empty queue, Len after pops, wrap-around, growth.
	_ = NewQueue[int]
}
