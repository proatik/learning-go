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
	return &Broker[T]{subs: make(map[int]chan T), buffer: max(buffer, 0)}
}

// Subscribe registers a new subscriber and returns its channel and a func
// that unsubscribes it (removing and closing the channel). The func is safe
// to call more than once. After Close, Subscribe returns a closed channel.
func (b *Broker[T]) Subscribe() (<-chan T, func()) {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan T, b.buffer)
	if b.closed {
		close(ch)
		return ch, func() {}
	}
	id := b.nextID
	b.nextID++
	b.subs[id] = ch
	return ch, func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		// Close may already have closed ch; only close what is still registered.
		if c, ok := b.subs[id]; ok {
			delete(b.subs, id)
			close(c)
		}
	}
}

// Publish offers v to every subscriber without blocking: a subscriber whose
// buffer is full misses v. It returns how many subscribers received v.
func (b *Broker[T]) Publish(v T) int {
	// Holding the lock while sending is what keeps an unsubscribe from
	// closing a channel mid-send. The sends never block, so it is short.
	b.mu.Lock()
	defer b.mu.Unlock()
	n := 0
	for _, ch := range b.subs {
		select {
		case ch <- v:
			n++
		default:
		}
	}
	return n
}

// Subscribers returns the number of current subscribers.
func (b *Broker[T]) Subscribers() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.subs)
}

// Close closes every subscriber channel. Later Publish calls deliver nothing.
// Close is safe to call more than once.
func (b *Broker[T]) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for id, ch := range b.subs {
		delete(b.subs, id)
		close(ch)
	}
}
