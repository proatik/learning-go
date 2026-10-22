# First Match

Check every item with a (possibly slow) predicate, all **concurrently**, and return the index of the **first** item, in input order, that matches.

```go
func FirstMatch[T any](items []T, pred func(T) bool) int
```

## Examples

| Call | Result |
|---|---|
| `FirstMatch([]int{1, 4, 6, 7}, isEven)` | `1` |
| `FirstMatch([]int{1, 3, 5}, isEven)` | `-1` |
| `FirstMatch([]string{"go", "rust", "zig"}, hasLen4)` (4 runes) | `1` |
| `FirstMatch(items, pred)` where `pred` sleeps 1s | the same answer, after **1s** in total |
| `FirstMatch([]int{}, pred)` / `FirstMatch(nil, pred)` | `-1`; `pred` is never called |

## Rules
- Call `pred` exactly once per item, each call in its own goroutine.
- Return the **smallest** index `i` with `pred(items[i]) == true`, or `-1` if none match.
- The answer must not depend on which goroutine finishes first: record each result at its index, decide after `Wait`.
