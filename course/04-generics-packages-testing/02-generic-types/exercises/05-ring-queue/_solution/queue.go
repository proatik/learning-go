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
	if q.n == len(q.buf) {
		q.grow()
	}
	q.buf[(q.head+q.n)%len(q.buf)] = v
	q.n++
}

// grow moves the elements, front first, into a buffer twice as long.
func (q *Queue[T]) grow() {
	buf := make([]T, max(4, 2*len(q.buf)))
	for i := range q.n {
		buf[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	q.buf, q.head = buf, 0
}

// Pop removes and returns the front value; ok is false if the queue is
// empty.
func (q *Queue[T]) Pop() (T, bool) {
	v, ok := q.Peek()
	if !ok {
		return v, false
	}
	var zero T
	q.buf[q.head] = zero
	q.head = (q.head + 1) % len(q.buf)
	q.n--
	return v, true
}

// Peek returns the front value without removing it; ok is false if the
// queue is empty.
func (q *Queue[T]) Peek() (T, bool) {
	if q.n == 0 {
		var zero T
		return zero, false
	}
	return q.buf[q.head], true
}

// Len returns the number of values in the queue.
func (q *Queue[T]) Len() int {
	return q.n
}

// Values returns a copy of the queue contents from front to back.
func (q *Queue[T]) Values() []T {
	out := make([]T, q.n)
	for i := range q.n {
		out[i] = q.buf[(q.head+i)%len(q.buf)]
	}
	return out
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
	return &Bounded[T]{limit: max(limit, 1)}
}

// Push adds v at the back, first dropping the oldest value if the queue is
// full. It shadows Queue.Push.
func (b *Bounded[T]) Push(v T) {
	if b.Len() == b.limit {
		b.Pop()
		b.dropped++
	}
	b.Queue.Push(v)
}

// Dropped returns how many values Push has dropped so far.
func (b *Bounded[T]) Dropped() int {
	return b.dropped
}
