# Priority Queue

Write a generic priority queue on your own **binary heap** (no `container/heap`), drain it with an iterator, and use it to find the top *k* values of any sequence.

```go
type Queue[T any] struct {
	items []T                // binary heap: items[i] is never less than items[(i-1)/2]
	less  func(a, b T) bool // a comes out before b
}

func New[T any](less func(a, b T) bool) *Queue[T]
func NewOrdered[T cmp.Ordered]() *Queue[T]

func (q *Queue[T]) Push(values ...T)
func (q *Queue[T]) Pop() (T, bool)
func (q *Queue[T]) Peek() (T, bool)
func (q *Queue[T]) Len() int
func (q *Queue[T]) Drain() iter.Seq[T]

func TopK[T cmp.Ordered](seq iter.Seq[T], k int) []T
```

## Examples

| Code | Result |
|---|---|
| `q := NewOrdered[int](); q.Push(3, 1, 2)` | `q.Peek()` → `1, true` |
| `slices.Collect(q.Drain())` | `[1 2 3]`, and `q.Len()` is `0` |
| `New(func(a, b task) bool { return a.pri > b.pri })` | highest `pri` comes out first |
| `for v := range q.Drain() { if v == 2 { break } }` on `1 2 3` | `q.Len()` → `1` (only `3` left) |
| `TopK(slices.Values([]int{5, 1, 9, 3}), 2)` | `[9 5]` |
| `TopK(seq, 0)` | `[]` (empty, non-nil) |

## Rules
- `items` must **always** be a valid heap: for every `i > 0`, `less(items[i], items[(i-1)/2])` is `false`. The tests inspect it.
- **NewOrdered** is a min-queue built with `New` and `cmp.Less[T]`.
- **Push** sifts each new value **up**; **Pop** moves the last item to the root and sifts it **down**. Both are O(log n).
- **Pop** / **Peek** on an empty queue return the zero value and `false`.
- **Drain** pops and yields until the queue is empty or `yield` returns `false`. A yielded item is gone; the rest stay queued.
- **TopK** returns at most `k` values, largest first. Use a min-queue of size `k`: when full, replace its smallest item only if the new value is bigger. `k <= 0` or an empty sequence → empty, non-nil slice.
