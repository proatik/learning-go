# Swap Through Pointers

A function that receives `int` values only gets copies. Receive **pointers** instead and change the caller's variables.

```go
func Swap(a, b *int)
func SwapIfGreater(a, b *int) bool
```

- `Swap` exchanges the values `a` and `b` point to.
- `SwapIfGreater` swaps only when `*a > *b` (so afterwards `*a <= *b`), and reports whether it swapped.

| Before | Call | After | Result |
|---|---|---|---|
| `x, y := 1, 2` | `Swap(&x, &y)` | `x == 2, y == 1` | — |
| `x := 7` | `Swap(&x, &x)` | `x == 7` | — |
| `x, y := 9, 4` | `SwapIfGreater(&x, &y)` | `x == 4, y == 9` | `true` |
| `x, y := 4, 9` | `SwapIfGreater(&x, &y)` | unchanged | `false` |
| `x, y := 5, 5` | `SwapIfGreater(&x, &y)` | unchanged | `false` |

## Rules
- Use a single tuple assignment: no temporary variable needed.
- Both pointers are always non-nil.
