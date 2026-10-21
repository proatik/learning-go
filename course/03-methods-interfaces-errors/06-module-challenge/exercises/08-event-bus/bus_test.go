package events_test

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"
	"testing"

	events "learngo/course/03-methods-interfaces-errors/06-module-challenge/exercises/08-event-bus"
)

// recorder is a Subscriber with a pointer receiver: only *recorder
// implements the interface.
type recorder struct {
	name string
	log  *[]string // shared log, to check the delivery order
}

func (r *recorder) Notify(e events.Event) error {
	*r.log = append(*r.log, fmt.Sprintf("%s<-%s:%v", r.name, e.Topic, e.Data))
	return nil
}

var errRejected = errors.New("rejected")

func newBus(t *testing.T) *events.Bus {
	t.Helper()
	b := events.NewBus()
	if b == nil {
		t.Fatalf("NewBus() = nil, want a bus")
	}
	return b
}

func subscribe(t *testing.T, b *events.Bus, topic string, s events.Subscriber, wantID int) {
	t.Helper()
	if got := b.Subscribe(topic, s); got != wantID {
		t.Fatalf("Subscribe(%q, ...) = %d, want %d", topic, got, wantID)
	}
}

func publish(t *testing.T, b *events.Bus, e events.Event, wantDelivered int) error {
	t.Helper()
	n, err := b.Publish(e)
	if n != wantDelivered {
		t.Errorf("Publish(%+v) delivered = %d, want %d (err = %v)", e, n, wantDelivered, err)
	}
	return err
}

func TestSubscriberFunc(t *testing.T) {
	var got events.Event
	f := events.SubscriberFunc(func(e events.Event) error {
		got = e
		return errRejected
	})
	var s events.Subscriber = f // a func type with a method satisfies the interface
	err := s.Notify(events.Event{Topic: "t", Data: 1})
	if err != errRejected || got.Topic != "t" || got.Data != 1 {
		t.Errorf("SubscriberFunc.Notify(t, 1) = %v (saw %+v), want %v (saw {t 1})", err, got, errRejected)
	}
}

func TestDeliveryOrderAndTopics(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "user", &recorder{"A", &log}, 1)
	subscribe(t, b, "*", &recorder{"W", &log}, 2)
	subscribe(t, b, "order", &recorder{"B", &log}, 3)
	subscribe(t, b, "user", &recorder{"C", &log}, 4)
	if got := b.Len(); got != 4 {
		t.Errorf("Len() = %d, want 4", got)
	}

	if err := publish(t, b, events.Event{Topic: "user", Data: "ada"}, 3); err != nil {
		t.Errorf("Publish(user) error = %v, want nil", err)
	}
	if err := publish(t, b, events.Event{Topic: "order", Data: 7}, 2); err != nil {
		t.Errorf("Publish(order) error = %v, want nil", err)
	}
	if err := publish(t, b, events.Event{Topic: "other"}, 1); err != nil {
		t.Errorf("Publish(other) error = %v, want nil", err)
	}
	want := []string{"A<-user:ada", "W<-user:ada", "C<-user:ada", "W<-order:7", "B<-order:7", "W<-other:<nil>"}
	if !slices.Equal(log, want) {
		t.Errorf("deliveries = %q, want %q", log, want)
	}
}

func TestSubscribeRejects(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "", &recorder{"A", &log}, 0)
	subscribe(t, b, "t", nil, 0)
	subscribe(t, b, "t", &recorder{"A", &log}, 1) // IDs are not used up by rejected calls
	if got := b.Len(); got != 1 {
		t.Errorf("Len() = %d, want 1", got)
	}
	n, err := b.Publish(events.Event{Topic: ""})
	if n != 0 || !errors.Is(err, events.ErrEmptyTopic) {
		t.Errorf("Publish(empty topic) = %d, %v, want 0, %v", n, err, events.ErrEmptyTopic)
	}
	if len(log) != 0 {
		t.Errorf("deliveries = %q, want none", log)
	}
}

func TestUnsubscribe(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "t", &recorder{"A", &log}, 1)
	subscribe(t, b, "t", &recorder{"B", &log}, 2)
	subscribe(t, b, "t", &recorder{"C", &log}, 3)
	if !b.Unsubscribe(2) {
		t.Errorf("Unsubscribe(2) = false, want true")
	}
	if b.Unsubscribe(2) {
		t.Errorf("Unsubscribe(2) twice = true, want false")
	}
	for _, id := range []int{0, -1, 99} {
		if b.Unsubscribe(id) {
			t.Errorf("Unsubscribe(%d) = true, want false", id)
		}
	}
	publish(t, b, events.Event{Topic: "t", Data: 1}, 2)
	subscribe(t, b, "t", &recorder{"D", &log}, 4) // IDs are never reused
	publish(t, b, events.Event{Topic: "t", Data: 2}, 3)
	want := []string{"A<-t:1", "C<-t:1", "A<-t:2", "C<-t:2", "D<-t:2"}
	if !slices.Equal(log, want) {
		t.Errorf("deliveries = %q, want %q", log, want)
	}
	if got := b.Len(); got != 3 {
		t.Errorf("Len() = %d, want 3", got)
	}
}

func TestErrorsAreCollected(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "t", &recorder{"A", &log}, 1)
	subscribe(t, b, "t", events.SubscriberFunc(func(events.Event) error { return errRejected }), 2)
	subscribe(t, b, "t", &recorder{"C", &log}, 3)
	subscribe(t, b, "*", events.SubscriberFunc(func(events.Event) error { return errors.New("full") }), 4)

	err := publish(t, b, events.Event{Topic: "t"}, 2)
	if err == nil {
		t.Fatalf("Publish error = nil, want the joined subscriber errors")
	}
	if want := "subscriber 2: rejected\nsubscriber 4: full"; err.Error() != want {
		t.Errorf("Publish error =\n%s\nwant\n%s", err.Error(), want)
	}
	if !errors.Is(err, errRejected) {
		t.Errorf("errors.Is(err, errRejected) = false, want true (wrap with %%w)")
	}
	if want := []string{"A<-t:<nil>", "C<-t:<nil>"}; !slices.Equal(log, want) {
		t.Errorf("deliveries = %q, want %q (a failing subscriber must not stop the others)", log, want)
	}
}

func TestPanicsAreIsolated(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "t", events.SubscriberFunc(func(events.Event) error { panic("boom") }), 1)
	subscribe(t, b, "t", &recorder{"B", &log}, 2)
	subscribe(t, b, "t", events.SubscriberFunc(func(events.Event) error { panic(errRejected) }), 3)
	subscribe(t, b, "t", &recorder{"D", &log}, 4)

	err := publish(t, b, events.Event{Topic: "t", Data: "x"}, 2)
	if want := []string{"B<-t:x", "D<-t:x"}; !slices.Equal(log, want) {
		t.Errorf("deliveries = %q, want %q (a panicking subscriber must not stop the others)", log, want)
	}
	if err == nil {
		t.Fatalf("Publish error = nil, want two *PanicError")
	}
	if want := "events: subscriber 1 panicked: boom\nevents: subscriber 3 panicked: rejected"; err.Error() != want {
		t.Errorf("Publish error =\n%s\nwant\n%s", err.Error(), want)
	}
	pe, ok := errors.AsType[*events.PanicError](err)
	if !ok || pe.ID != 1 || pe.Value != "boom" {
		t.Errorf("errors.AsType[*PanicError](err) = %+v, %v, want {ID:1 Value:boom}, true", pe, ok)
	}
	// A panic value that is an error is reachable through Unwrap.
	if !errors.Is(err, errRejected) {
		t.Errorf("errors.Is(err, errRejected) = false, want true (PanicError.Unwrap)")
	}
	if got := (&events.PanicError{ID: 9, Value: 42}).Unwrap(); got != nil {
		t.Errorf("PanicError{Value: 42}.Unwrap() = %v, want nil", got)
	}

	// The bus keeps working afterwards.
	log = nil
	b.Unsubscribe(1)
	b.Unsubscribe(3)
	if err := publish(t, b, events.Event{Topic: "t", Data: "y"}, 2); err != nil {
		t.Errorf("Publish after removing the panicking subscribers error = %v, want nil", err)
	}
}

func TestTypedNilSubscriber(t *testing.T) {
	b := newBus(t)
	var r *recorder // nil pointer, but a non-nil Subscriber interface value
	subscribe(t, b, "t", r, 1)
	err := publish(t, b, events.Event{Topic: "t"}, 0)
	if _, ok := errors.AsType[*events.PanicError](err); !ok {
		t.Fatalf("Publish to a typed-nil subscriber error = %v, want a *PanicError", err)
	}
	// The recovered value is a runtime error (nil pointer dereference).
	if re, ok := errors.AsType[runtime.Error](err); !ok || !strings.Contains(re.Error(), "nil pointer") {
		t.Errorf("errors.AsType[runtime.Error](err) = %v, %v, want the nil-pointer runtime error, true", re, ok)
	}
}

func TestUnsubscribeDuringPublish(t *testing.T) {
	b := newBus(t)
	var log []string
	subscribe(t, b, "t", events.SubscriberFunc(func(events.Event) error {
		b.Unsubscribe(1) // unsubscribe itself
		b.Unsubscribe(2) // and the next one
		return nil
	}), 1)
	subscribe(t, b, "t", &recorder{"B", &log}, 2)
	// Publish works on the subscriptions that existed when it started.
	publish(t, b, events.Event{Topic: "t", Data: 1}, 2)
	publish(t, b, events.Event{Topic: "t", Data: 2}, 0)
	if want := []string{"B<-t:1"}; !slices.Equal(log, want) {
		t.Errorf("deliveries = %q, want %q", log, want)
	}
}
