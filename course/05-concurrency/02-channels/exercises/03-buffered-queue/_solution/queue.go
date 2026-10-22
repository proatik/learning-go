package queue

// Queue is a fixed-size FIFO queue of strings backed by a buffered channel.
// It is meant for use by a single goroutine.
type Queue struct {
	ch chan string
}

// New returns an empty queue that holds up to capacity values.
func New(capacity int) *Queue {
	return &Queue{ch: make(chan string, capacity)}
}

// Push appends s and reports true, or reports false without blocking if the
// queue is full.
func (q *Queue) Push(s string) bool {
	if len(q.ch) == cap(q.ch) {
		return false
	}
	q.ch <- s
	return true
}

// Pop removes and returns the oldest value, or returns "", false without
// blocking if the queue is empty.
func (q *Queue) Pop() (string, bool) {
	if len(q.ch) == 0 {
		return "", false
	}
	return <-q.ch, true
}

// Len returns the number of queued values.
func (q *Queue) Len() int {
	return len(q.ch)
}

// Cap returns the maximum number of values the queue can hold.
func (q *Queue) Cap() int {
	return cap(q.ch)
}
