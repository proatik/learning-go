# Buffered Queue

Use a **buffered channel** as a fixed-size FIFO queue whose `Push` and `Pop` never block: they check `len` and `cap` first and report failure instead.

```go
type Queue struct {
	ch chan string
}

func New(capacity int) *Queue

func (q *Queue) Push(s string) bool   // false if the queue is full
func (q *Queue) Pop() (string, bool)  // "", false if the queue is empty
func (q *Queue) Len() int
func (q *Queue) Cap() int
```

## Examples

| Code (`q := New(2)`) | Result |
|---|---|
| `q.Len(), q.Cap()` | `0, 2` |
| `q.Push("a"); q.Push("b")` | `true`, `true` |
| `q.Push("c")` | `false` (full, nothing changes) |
| `q.Len()` | `2` |
| `q.Pop()` | `"a", true` (first in, first out) |
| `q.Pop()`, `q.Pop()` | `"b", true`, then `"", false` |
| `New(0).Push("x")` | `false`: an unbuffered channel has no room at all |

## Rules
- `New` creates the channel with `make(chan string, capacity)`. `capacity` is never negative in the tests.
- `Push` adds `s` and returns `true`, or returns `false` **without blocking** when `Len() == Cap()`.
- `Pop` removes the oldest value and returns it with `true`, or returns `"", false` **without blocking** when the queue is empty. An empty string is a valid value.
- `Len` is the number of queued values, `Cap` the capacity.
- The queue is used by **one goroutine** only: with several, the check and the send could be interleaved by another goroutine (lesson 3's `select` with `default` solves that).
