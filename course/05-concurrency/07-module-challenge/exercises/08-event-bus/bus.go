package bus

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrClosed is returned by Publish after Close.
var ErrClosed = errors.New("bus: closed")

// Bus delivers messages published to a topic to every subscriber of that
// topic. Create one with New.
type Bus struct {
	mu        sync.Mutex
	queueSize int
	subs      map[string]map[*Subscription]bool // topic → subscriptions
	closed    bool
}

// Subscription is one handler subscribed to one topic.
type Subscription struct {
	bus     *Bus
	topic   string
	queue   chan string   // messages waiting for the handler
	done    chan struct{} // closed when the handler goroutine exits
	dropped atomic.Int64
	once    sync.Once
}

// New returns a bus whose subscriptions each queue up to queueSize messages.
func New(queueSize int) *Bus {
	// TODO
	return &Bus{}
}

// Subscribe calls handler, on its own goroutine, for every message published
// to topic from now on.
func (b *Bus) Subscribe(topic string, handler func(msg string)) *Subscription {
	// TODO
	return &Subscription{}
}

// Publish queues msg for every subscriber of topic without blocking,
// dropping a subscriber's oldest queued message if its queue is full.
func (b *Bus) Publish(topic, msg string) error {
	// TODO
	return nil
}

// Close unsubscribes everyone and waits until every queued message has been
// handled. It is safe to call more than once.
func (b *Bus) Close() {
	// TODO
}

// Unsubscribe stops new messages, waits until the already queued ones have
// been handled, and returns. It is safe to call more than once.
func (s *Subscription) Unsubscribe() {
	// TODO
}

// Dropped returns how many messages were dropped because the queue was full.
func (s *Subscription) Dropped() int {
	// TODO
	return 0
}
