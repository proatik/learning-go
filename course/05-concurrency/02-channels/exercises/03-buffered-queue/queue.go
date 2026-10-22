package queue

// Queue is a fixed-size FIFO queue of strings backed by a buffered channel.
// It is meant for use by a single goroutine.
type Queue struct {
	ch chan string
}

// New returns an empty queue that holds up to capacity values.
func New(capacity int) *Queue {
	// TODO
	return &Queue{}
}

// Push appends s and reports true, or reports false without blocking if the
// queue is full.
func (q *Queue) Push(s string) bool {
	// TODO
	return false
}

// Pop removes and returns the oldest value, or returns "", false without
// blocking if the queue is empty.
func (q *Queue) Pop() (string, bool) {
	// TODO
	return "", false
}

// Len returns the number of queued values.
func (q *Queue) Len() int {
	// TODO
	return 0
}

// Cap returns the maximum number of values the queue can hold.
func (q *Queue) Cap() int {
	// TODO
	return 0
}
