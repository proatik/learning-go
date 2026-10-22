package pubsub

import "sync"

// Broker delivers every published value to every current subscriber. It is
// safe for concurrent use. Create one with New.
type Broker[T any] struct {
	mu     sync.Mutex
	subs   map[int]chan T // subscriber id → its channel
	nextID int
	buffer int
	closed bool
}

// New returns a broker whose subscriber channels have the given buffer size
// (buffer < 0 means 0).
func New[T any](buffer int) *Broker[T] {
	// TODO
	return &Broker[T]{}
}

// Subscribe registers a new subscriber and returns its channel and a func
// that unsubscribes it (removing and closing the channel). The func is safe
// to call more than once. After Close, Subscribe returns a closed channel.
func (b *Broker[T]) Subscribe() (<-chan T, func()) {
	// TODO
	return nil, func() {}
}

// Publish offers v to every subscriber without blocking: a subscriber whose
// buffer is full misses v. It returns how many subscribers received v.
func (b *Broker[T]) Publish(v T) int {
	// TODO
	return 0
}

// Subscribers returns the number of current subscribers.
func (b *Broker[T]) Subscribers() int {
	// TODO
	return 0
}

// Close closes every subscriber channel. Later Publish calls deliver nothing.
// Close is safe to call more than once.
func (b *Broker[T]) Close() {
	// TODO
}
