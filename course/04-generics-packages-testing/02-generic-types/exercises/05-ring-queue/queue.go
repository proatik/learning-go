package queue

// Queue is a first-in, first-out queue stored in a growable ring buffer. The
// zero value is an empty queue.
type Queue[T any] struct {
	buf  []T // ring storage; len(buf) is the capacity
	head int // index of the front element in buf
	n    int // number of stored elements
}

// Push adds v at the back, growing the buffer to max(4, 2*len(buf)) when it
// is full.
func (q *Queue[T]) Push(v T) {
	// TODO
}

// Pop removes and returns the front value; ok is false if the queue is
// empty.
func (q *Queue[T]) Pop() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Peek returns the front value without removing it; ok is false if the
// queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	// TODO
	var zero T
	return zero, false
}

// Len returns the number of values in the queue.
func (q *Queue[T]) Len() int {
	// TODO
	return 0
}

// Values returns a copy of the queue contents from front to back.
func (q *Queue[T]) Values() []T {
	// TODO
	return nil
}

// Bounded is a Queue that holds at most limit values: pushing onto a full
// Bounded drops the oldest value. Pop, Peek, Len and Values are promoted
// from the embedded Queue.
type Bounded[T any] struct {
	Queue[T]
	limit   int
	dropped int
}

// NewBounded returns an empty Bounded queue; a limit below 1 is treated as 1.
func NewBounded[T any](limit int) *Bounded[T] {
	// TODO
	return &Bounded[T]{}
}

// Push adds v at the back, first dropping the oldest value if the queue is
// full. It shadows Queue.Push.
func (b *Bounded[T]) Push(v T) {
	// TODO
}

// Dropped returns how many values Push has dropped so far.
func (b *Bounded[T]) Dropped() int {
	// TODO
	return 0
}
