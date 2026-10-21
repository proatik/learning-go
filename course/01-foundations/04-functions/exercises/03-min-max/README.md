# Min and Max in One Pass

Return the smallest and the largest of one or more numbers:

```go
func MinMax(first int, rest ...int) (lo, hi int)
```

The separate `first` parameter makes "at least one number" part of the signature: `MinMax()` doesn't compile, so there is no empty case to handle.

| Call | `lo` | `hi` |
|---|---|---|
| `MinMax(7)` | `7` | `7` |
| `MinMax(3, 9, -2, 5)` | `-2` | `9` |
| `MinMax(4, 4, 4)` | `4` | `4` |
| `MinMax(5, xs...)` with `xs := []int{1, 8}` | `1` | `8` |

## Rules
- Start with `lo, hi = first, first`, then loop over `rest` once.
- Use named results and a final `return`.
- Write the comparisons yourself with `if` (the built-in `min`/`max` functions exist, but practise the loop here).
