# Pub/Sub Broker

A broker that **broadcasts**: every published value goes to every current subscriber, each on its own buffered channel. Slow subscribers never block the publisher; they just miss values.

```go
type Broker[T any] struct {
	mu     sync.Mutex
	subs   map[int]chan T // subscriber id → its channel
	nextID int
	buffer int
	closed bool
}

func New[T any](buffer int) *Broker[T]
func (b *Broker[T]) Subscribe() (<-chan T, func())
func (b *Broker[T]) Publish(v T) int
func (b *Broker[T]) Subscribers() int
func (b *Broker[T]) Close()
```

```text
                 ┌──▶ [buf] ──▶ subscriber 1
Publish(v) ──────┼──▶ [buf] ──▶ subscriber 2
                 └──▶ [full] ✗  subscriber 3 misses v
```

## Examples

| Code (buffer 10 unless noted) | Result |
|---|---|
| `Publish(99)` with no subscribers | `0` |
| two subscribers, `Publish(1)`, `Publish(2)` | each returns `2`; both channels hold `1 2` |
| `ch, unsub := Subscribe(); unsub()` | `ch` is closed (after any buffered values); calling `unsub()` again does nothing |
| buffer 2, a subscriber that never reads, `Publish(1..4)` | it keeps `1 2`; `Publish(3)` and `Publish(4)` skip it |
| buffer 0, a subscriber blocked in `<-ch` | `Publish(v)` returns `1` (it delivers only to a waiting receiver) |
| `Close()` | every subscriber channel is closed; `Publish` returns `0`; `Subscribers()` is `0` |
| `Subscribe()` after `Close()` | a closed channel and a func that does nothing |

## Rules
- `New` clamps a negative buffer to `0` and makes the map.
- `Publish` **never blocks**: for each subscriber, a `select` with a `default` case either sends or skips. It returns the number of subscribers that got the value.
- The broker owns every subscriber channel, so the broker is the one that closes it: in `unsubscribe` or in `Close`, **exactly once**. The unsubscribe func and `Close` must both be safe to call more than once and in any order.
- Every method is safe for concurrent use. Keep the mutex held during `Publish`'s sends: otherwise an unsubscribe could close a channel in the middle of a send, and sending on a closed channel panics.
