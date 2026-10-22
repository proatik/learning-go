# Stack Iterators

Add iterator methods to a generic stack so callers can `range` over it.

```go
type Stack[T any] struct {
	items []T // items[len(items)-1] is the top
}

func (s *Stack[T]) Push(v T)          // provided
func (s *Stack[T]) Pop() (T, bool)    // provided
func (s *Stack[T]) Len() int          // provided

func (s *Stack[T]) All() iter.Seq[T]      // top → bottom (the order Pop would return them)
func (s *Stack[T]) Backward() iter.Seq[T] // bottom → top (the order they were pushed)
func (s *Stack[T]) Drain() iter.Seq[T]    // pops and yields, top → bottom
```

## Examples

Push `"a"`, `"b"`, `"c"` onto `s`:

| Code | Values | `s.Len()` afterwards |
|---|---|---|
| `slices.Collect(s.All())` | `[c b a]` | 3 |
| `slices.Collect(s.Backward())` | `[a b c]` | 3 |
| `slices.Collect(s.Drain())` | `[c b a]` | 0 |
| `for v := range s.Drain() { if v == "b" { break } }` | `c`, `b` | 1 (`"a"` is still on the stack) |

## Rules
- `All` and `Backward` must not change the stack. Ranging over them twice gives the same values.
- `Drain` pops one item, yields it, and repeats until the stack is empty or `yield` returns `false`. Items that were never yielded stay on the stack.
- The zero value `Stack[T]{}` is an empty stack: every iterator yields nothing.
- Stop as soon as `yield` returns `false`.
