# Sign, Abs and Clamp

Implement:

```go
func Sign(n int) int          // -1, 0 or 1
func Abs(n int) int           // distance from zero
func Clamp(x, lo, hi int) int // x limited to the range [lo, hi]
```

## Examples
| Call | Result |
|---|---|
| `Sign(-7)` | `-1` |
| `Sign(0)` | `0` |
| `Sign(42)` | `1` |
| `Abs(-5)` | `5` |
| `Abs(3)` | `3` |
| `Clamp(15, 0, 10)` | `10` |
| `Clamp(-3, 0, 10)` | `0` |
| `Clamp(4, 0, 10)` | `4` |

## Rules
- `Sign`: use `if` / `else if` / `else`.
- `Abs`: a single `if` with an early `return` is enough.
- `Clamp`: you may assume `lo <= hi`. Solve it with the builtins `min` and `max` (no `if`).
