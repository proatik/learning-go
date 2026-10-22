# Ring-Buffer Queue

A first-in, first-out `Queue[T]` stored in a **ring buffer** that grows when full, so `Pop` never shifts elements. Then **embed** it in a `Bounded[T]` queue that drops its oldest value when full.

```go
type Queue[T any] struct {
	buf  []T // ring storage; len(buf) is the capacity
	head int // index of the front element in buf
	n    int // number of stored elements
}

func (q *Queue[T]) Push(v T)
func (q *Queue[T]) Pop() (T, bool)
func (q *Queue[T]) Peek() (T, bool)
func (q *Queue[T]) Len() int
func (q *Queue[T]) Values() []T

type Bounded[T any] struct {
	Queue[T]
	limit   int
	dropped int
}

func NewBounded[T any](limit int) *Bounded[T]
func (b *Bounded[T]) Push(v T)
func (b *Bounded[T]) Dropped() int
```

## Examples

| Code | Result |
|---|---|
| `var q Queue[int]; q.Push(1); q.Push(2); q.Push(3)` | `q.Values()` = `[1 2 3]` |
| `q.Pop()` | `1, true` |
| `q.Peek()` | `2, true` (queue unchanged) |
| `var e Queue[string]; e.Pop()` | `"", false` |
| `b := NewBounded[int](2); b.Push(1); b.Push(2); b.Push(3)` | `b.Values()` = `[2 3]`, `b.Dropped()` = `1` |
| `b.Pop()` | `2, true` (promoted from `Queue[int]`) |

## Rules
- The zero `Queue[T]` is empty and usable.
- **Ring buffer:** the element after index `i` is at `(i+1) % len(q.buf)`. `Push` writes at `(head+n) % len(buf)`; `Pop` reads at `head` and advances it. Never shift or re-slice the elements.
- **Growth:** when `Push` finds the buffer full (`n == len(buf)`), allocate a new buffer of length `max(4, 2*len(buf))`, copy the elements into it **front first starting at index 0**, and set `head = 0`. The buffer never grows otherwise (tests check `len(q.buf)`).
- `Pop` must set the vacated slot to the zero value of `T`, so the queue doesn't keep popped pointers alive.
- `Pop` / `Peek` on an empty queue return the zero value of `T` and `false`.
- `Values` returns a **new** slice, front to back (empty or nil when empty).
- `NewBounded` treats `limit < 1` as `1`.
- `Bounded.Push` first pops the oldest value if `Len() == limit` (counting it in `Dropped`), then pushes `v` onto the embedded queue. All other methods are promoted from `Queue[T]`.
