package events

import (
	"errors"
	"fmt"
	"slices"
)

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
func (f SubscriberFunc) Notify(e Event) error { return f(e) }

// PanicError records a subscriber that panicked.
type PanicError struct {
	ID    int
	Value any
}

// Error returns "events: subscriber ID panicked: VALUE".
func (e *PanicError) Error() string {
	return fmt.Sprintf("events: subscriber %d panicked: %v", e.ID, e.Value)
}

// Unwrap returns the panic value if it is an error, otherwise nil.
func (e *PanicError) Unwrap() error {
	err, _ := e.Value.(error)
	return err
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
func NewBus() *Bus { return &Bus{} }

// Subscribe registers s for topic ("*" means every topic) and returns its
// ID, starting at 1. An empty topic or nil s registers nothing and returns 0.
func (b *Bus) Subscribe(topic string, s Subscriber) int {
	if topic == "" || s == nil {
		return 0
	}
	b.nextID++
	b.subs = append(b.subs, subscription{id: b.nextID, topic: topic, sub: s})
	return b.nextID
}

// Unsubscribe removes the subscription with the given ID and reports
// whether it existed.
func (b *Bus) Unsubscribe(id int) bool {
	before := len(b.subs)
	b.subs = slices.DeleteFunc(b.subs, func(s subscription) bool { return s.id == id })
	return len(b.subs) < before
}

// Len returns the number of subscriptions.
func (b *Bus) Len() int { return len(b.subs) }

// notify calls one subscriber, turning a panic into a *PanicError.
func notify(s subscription, e Event) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &PanicError{ID: s.id, Value: r}
		}
	}()
	if err := s.sub.Notify(e); err != nil {
		return fmt.Errorf("subscriber %d: %w", s.id, err)
	}
	return nil
}

// Publish delivers e to every matching subscriber in subscription order.
// It returns how many succeeded and all failures joined with errors.Join.
func (b *Bus) Publish(e Event) (delivered int, err error) {
	if e.Topic == "" {
		return 0, ErrEmptyTopic
	}
	var errs []error
	for _, s := range slices.Clone(b.subs) {
		if s.topic != e.Topic && s.topic != "*" {
			continue
		}
		if err := notify(s, e); err != nil {
			errs = append(errs, err)
			continue
		}
		delivered++
	}
	return delivered, errors.Join(errs...)
}
