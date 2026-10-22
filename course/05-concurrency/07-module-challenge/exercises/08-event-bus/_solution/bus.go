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
	return &Bus{queueSize: max(queueSize, 1), subs: map[string]map[*Subscription]bool{}}
}

// Subscribe calls handler, on its own goroutine, for every message published
// to topic from now on.
func (b *Bus) Subscribe(topic string, handler func(msg string)) *Subscription {
	s := &Subscription{
		bus:   b,
		topic: topic,
		queue: make(chan string, b.queueSize),
		done:  make(chan struct{}),
	}
	go func() {
		defer close(s.done)
		for msg := range s.queue {
			handler(msg)
		}
	}()

	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		s.once.Do(func() { close(s.queue) })
		return s
	}
	if b.subs[topic] == nil {
		b.subs[topic] = map[*Subscription]bool{}
	}
	b.subs[topic][s] = true
	return s
}

// Publish queues msg for every subscriber of topic without blocking,
// dropping a subscriber's oldest queued message if its queue is full.
func (b *Bus) Publish(topic, msg string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return ErrClosed
	}
	for s := range b.subs[topic] {
		s.enqueue(msg)
	}
	return nil
}

// enqueue adds msg to the queue, dropping the oldest message while it is
// full. The handler may take a message at the same time, so loop.
func (s *Subscription) enqueue(msg string) {
	for {
		select {
		case s.queue <- msg:
			return
		default:
		}
		select {
		case <-s.queue:
			s.dropped.Add(1)
		default:
		}
	}
}

// Close unsubscribes everyone and waits until every queued message has been
// handled. It is safe to call more than once.
func (b *Bus) Close() {
	b.mu.Lock()
	b.closed = true
	var all []*Subscription
	for _, set := range b.subs {
		for s := range set {
			all = append(all, s)
		}
	}
	b.mu.Unlock()
	for _, s := range all {
		s.Unsubscribe()
	}
}

// Unsubscribe stops new messages, waits until the already queued ones have
// been handled, and returns. It is safe to call more than once.
func (s *Subscription) Unsubscribe() {
	b := s.bus
	b.mu.Lock()
	delete(b.subs[s.topic], s)
	if len(b.subs[s.topic]) == 0 {
		delete(b.subs, s.topic)
	}
	s.once.Do(func() { close(s.queue) }) // under mu: Publish can't be sending
	b.mu.Unlock()
	<-s.done
}

// Dropped returns how many messages were dropped because the queue was full.
func (s *Subscription) Dropped() int {
	return int(s.dropped.Load())
}
