# Event Bus

An in-process **event bus**: handlers subscribe to topics, publishers send messages to topics. Every subscription has its own goroutine and a bounded queue, so a **slow handler never blocks a publisher**: when its queue is full, the **oldest** queued message is dropped (and counted) to make room.

```go
var ErrClosed = errors.New("bus: closed")

type Bus struct {
	mu        sync.Mutex
	queueSize int
	subs      map[string]map[*Subscription]bool // topic → subscriptions
	closed    bool
}

type Subscription struct {
	bus     *Bus
	topic   string
	queue   chan string   // messages waiting for the handler
	done    chan struct{} // closed when the handler goroutine exits
	dropped atomic.Int64
	once    sync.Once
}

func New(queueSize int) *Bus

func (b *Bus) Subscribe(topic string, handler func(msg string)) *Subscription
func (b *Bus) Publish(topic, msg string) error
func (b *Bus) Close()

func (s *Subscription) Unsubscribe()
func (s *Subscription) Dropped() int
```

## Examples

| Code (`b := New(2)`; one subscriber on `"news"` whose handler takes 1s per message) | Result |
|---|---|
| `b.Publish("news", "m1")` | returns at once; the handler starts on `m1` |
| `b.Publish` `"m2"`, `"m3"`, `"m4"`, `"m5"` while it's busy | each returns at once; queue holds 2: `m4` pushes out `m2`, `m5` pushes out `m3` |
| `s.Dropped()` | `2` |
| `b.Close()` | waits for the handler to finish `m1`, `m4`, `m5` (3s), then returns |
| `b.Publish("news", "m6")` after `Close` | `ErrClosed` |
| `b.Publish("sport", "x")` with no subscribers | `nil` (nobody gets it) |

## Rules
- `New(queueSize)`: every subscription's queue holds `queueSize` messages (`< 1` counts as `1`).
- `Subscribe` registers `handler` for `topic` and starts its goroutine, which calls `handler` for each queued message, one at a time, **in publish order**. Subscribing to a closed bus returns a subscription whose handler is never called (its goroutine still starts and exits).
- `Publish` delivers `msg` to every current subscriber of `topic` **without ever blocking** on a handler: if a queue is full, take out the oldest message, count it in that subscription's `Dropped`, and queue the new one. It returns `ErrClosed` after `Close`, otherwise `nil`.
- `Unsubscribe` removes the subscription so no new messages reach it, lets the handler finish the messages **already queued**, and returns after its goroutine has exited. Calling it again (or after `Close`) does nothing. (Calling it from inside its own handler would wait for itself forever; the tests never do.)
- `Close` marks the bus closed, unsubscribes everyone the same way (queued messages are still handled) and returns after every handler goroutine has exited. Calling it again does nothing.
- Never send on a closed channel: queues are closed only while holding `mu`, and `Publish` sends only while holding `mu`.
- Safe for concurrent use.
