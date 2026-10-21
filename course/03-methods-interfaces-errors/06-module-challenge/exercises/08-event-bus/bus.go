package events

import "errors"

// ErrEmptyTopic is returned by Publish for an event without a topic.
var ErrEmptyTopic = errors.New("events: empty topic")

// Event is a message published on a topic.
type Event struct {
	Topic string
	Data  any
}

// Subscriber receives events.
type Subscriber interface {
	Notify(e Event) error
}

// SubscriberFunc adapts an ordinary function to the Subscriber interface.
type SubscriberFunc func(e Event) error

// Notify calls f(e).
func (f SubscriberFunc) Notify(e Event) error {
	// TODO
	return nil
}

// PanicError records a subscriber that panicked.
type PanicError struct {
	ID    int
	Value any
}

// Error returns "events: subscriber ID panicked: VALUE".
func (e *PanicError) Error() string {
	// TODO
	return ""
}

// Unwrap returns the panic value if it is an error, otherwise nil.
func (e *PanicError) Unwrap() error {
	// TODO
	return nil
}

type subscription struct {
	id    int
	topic string
	sub   Subscriber
}

// Bus delivers published events to subscribers. Create it with NewBus.
type Bus struct {
	subs   []subscription
	nextID int
}

// NewBus returns an empty bus.
func NewBus() *Bus {
	// TODO
	return nil
}

// Subscribe registers s for topic ("*" means every topic) and returns its
// ID, starting at 1. An empty topic or nil s registers nothing and returns 0.
func (b *Bus) Subscribe(topic string, s Subscriber) int {
	// TODO
	return 0
}

// Unsubscribe removes the subscription with the given ID and reports
// whether it existed.
func (b *Bus) Unsubscribe(id int) bool {
	// TODO
	return false
}

// Len returns the number of subscriptions.
func (b *Bus) Len() int {
	// TODO
	return 0
}

// Publish delivers e to every matching subscriber in subscription order.
// It returns how many succeeded and all failures joined with errors.Join.
func (b *Bus) Publish(e Event) (delivered int, err error) {
	// TODO
	return 0, nil
}
