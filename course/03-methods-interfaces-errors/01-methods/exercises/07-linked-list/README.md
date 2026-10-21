# Nil-Safe Linked List

Module 2's linked list, rebuilt with methods. A **nil `*List` is the empty list**, so every method must work on a nil receiver.

```go
type List struct {
	Val  int
	Next *List
}

func FromSlice(vals []int) *List

func (l *List) Len() int
func (l *List) Sum() int
func (l *List) Contains(v int) bool
func (l *List) Values() []int
func (l *List) String() string
func (l *List) Push(v int) *List      // prepend, returns the new head
func (l *List) Append(v int) *List    // add at the end, returns the head
func (l *List) Reverse() *List        // in place, returns the new head
func (l *List) RemoveAll(v int) *List // returns the new head
```

## Examples

| Code | Result |
|---|---|
| `var l *List; l.Len()` | `0` (no panic) |
| `FromSlice([]int{1, 2, 3}).Values()` | `[1 2 3]` |
| `FromSlice(nil)` | `nil` |
| `FromSlice([]int{1, 2, 3}).String()` | `"[1 -> 2 -> 3]"` |
| `l.String()` with `l == nil` | `"[]"` |
| `l.Push(0)` on `[1 2]` | `[0 1 2]` |
| `l.Append(5)` on nil | `[5]` |
| `l.Reverse()` on `[1 2 3]` | `[3 2 1]` |
| `l.RemoveAll(2)` on `[2 1 2 3 2]` | `[1 3]` |

## Rules
- Every method must handle `l == nil` without panicking.
- `Values` returns `nil` for the empty list.
- `Push`, `Append`, `Reverse` and `RemoveAll` reuse the existing nodes (no copying) and return the head to use afterwards: `l = l.Push(4)`.
- `Reverse` must be iterative or recursive in place; don't build a new list.
