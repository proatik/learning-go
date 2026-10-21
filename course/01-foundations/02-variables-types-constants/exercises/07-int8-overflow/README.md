# Detect Integer Overflow

An `int8` holds −128 … 127. When a sum goes past that, Go silently **wraps around**: `int8(127) + 1 == -128`.

Implement:

```go
func AddInt8(a, b int8) (sum int8, overflow bool) // wrapped sum + did it overflow?
func SaturatingAdd(a, b int8) int8                // clamp to [MinInt8, MaxInt8] instead of wrapping
```

| Call | `AddInt8` | `SaturatingAdd` |
|---|---|---|
| `(1, 2)` | `3, false` | `3` |
| `(100, 27)` | `127, false` | `127` |
| `(100, 28)` | `-128, true` | `127` |
| `(127, 127)` | `-2, true` | `127` |
| `(-100, -29)` | `127, true` | `-128` |
| `(-128, -1)` | `127, true` | `-128` |
| `(-128, 127)` | `-1, false` | `-1` |

The result names `(sum int8, overflow bool)` are just documentation here; `return x, y` works as usual.

## Rules
- Use the limits from package `math`: `math.MaxInt8`, `math.MinInt8`.
- Tip: an `int16` can hold any sum of two `int8`s.
- The builtins `min(x, y)` and `max(x, y)` return the smaller / larger of their arguments (all of the same type).
- No `if` is needed in `SaturatingAdd`, though you may use one.
