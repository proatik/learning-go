# Event Bus

Build a publish/subscribe bus (the observer pattern). Subscribers are anything with a `Notify` method — including plain functions via an adapter type. A subscriber that fails or **panics** must not break delivery to the others: recover the panic, turn it into an error, and return every failure with `errors.Join`.

```go
var ErrEmptyTopic = errors.New("events: empty topic")

type Event struct {
	Topic string
	Data  any
}

type Subscriber interface {
	Notify(e Event) error
}

type SubscriberFunc func(e Event) error
func (f SubscriberFunc) Notify(e Event) error

type PanicError struct {
	ID    int // subscription ID
	Value any // what recover() returned
}
func (e *PanicError) Error() string
func (e *PanicError) Unwrap() error

type subscription struct { id int; topic string; sub Subscriber }
type Bus struct { subs []subscription; nextID int }

func NewBus() *Bus
func (b *Bus) Subscribe(topic string, s Subscriber) int
func (b *Bus) Unsubscribe(id int) bool
func (b *Bus) Len() int
func (b *Bus) Publish(e Event) (delivered int, err error)
```

## Examples

| Code | Result |
|---|---|
| `b.Subscribe("user", &rec)` | `1` |
| `b.Subscribe("*", SubscriberFunc(f))` | `2` — `"*"` receives every topic |
| `b.Publish(Event{Topic: "user"})` | delivers to 1 then 2: `2, nil` |
| subscriber 2 returns `errRejected` | `1, "subscriber 2: rejected"` |
| subscriber 1 does `panic("boom")` | `*PanicError` — `"events: subscriber 1 panicked: boom"` |
| `b.Publish(Event{})` | `0, ErrEmptyTopic` |

## Rules
- **SubscriberFunc.Notify** just calls the function — the same trick as `http.HandlerFunc`.
- **Subscribe** returns IDs `1, 2, 3, …` (never reused). An empty topic or a `nil` subscriber registers nothing and returns `0` without using up an ID.
- **Unsubscribe** removes that subscription and reports whether it existed. **Len** counts subscriptions.
- **Publish**:
  - empty topic → `0, ErrEmptyTopic`
  - deliver to every subscription whose topic equals `e.Topic` **or** is `"*"`, in **subscription order**
  - work on a snapshot (`slices.Clone`): subscriptions changed during `Publish` affect only later publishes
  - a returned error becomes `fmt.Errorf("subscriber %d: %w", id, err)`
  - a panic is recovered and becomes `&PanicError{ID: id, Value: r}`; delivery continues with the next subscriber
  - `delivered` counts subscribers that returned `nil` without panicking; `err` is `errors.Join` of all failures in order (`nil` if none)
- **PanicError.Error**: `"events: subscriber %d panicked: %v"`. **Unwrap** returns `Value` if it is an `error` (so `errors.Is` can see it), otherwise `nil`.
- Watch out: a nil `*T` stored in a `Subscriber` is **not** a nil interface, so `Subscribe` accepts it — and its `Notify` may panic at publish time. That panic must be isolated like any other.
