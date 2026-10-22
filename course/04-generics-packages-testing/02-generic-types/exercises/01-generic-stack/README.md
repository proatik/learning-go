# Generic Stack

Turn Module 3's `int` stack into a `Stack[T]` that works for any element type, and drain any container through a generic interface.

```go
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(v T)
func (s *Stack[T]) Pop() (T, bool)
func (s *Stack[T]) Peek() (T, bool)
func (s *Stack[T]) Len() int

type Container[T any] interface {
	Push(T)
	Pop() (T, bool)
	Len() int
}

func Drain[T any](c Container[T]) []T
```

## Examples

| Code | Result |
|---|---|
| `var s Stack[string]; s.Push("a"); s.Push("b")` | `s.Len()` = `2` |
| `s.Peek()` | `"b", true` (stack unchanged) |
| `s.Pop()` | `"b", true`, then `s.Len()` = `1` |
| `var e Stack[int]; e.Pop()` | `0, false` |
| `var p Stack[*int]; p.Pop()` | `nil, false` |
| `Drain[int](s)` with `s` holding `1 2 3` (3 on top) | `[3 2 1]`, and `s.Len()` = `0` afterwards |

## Rules
- The zero value (`var s Stack[T]`) is an empty, usable stack.
- `Pop` / `Peek` on an empty stack return the **zero value of `T`** and `false`.
- `*Stack[T]` must satisfy `Container[T]`.
- `Drain` pops every value from `c` and returns them **in pop order**. It must use only the `Container[T]` methods (the tests pass other container types too).
- `Drain` on an empty container returns an empty or nil slice.
