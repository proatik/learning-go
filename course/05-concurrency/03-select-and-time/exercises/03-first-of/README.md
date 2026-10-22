# First Of Two

Wait for whichever of two channels delivers first, and say which one it was. Give up after `timeout`.

```go
func FirstOf[T any](a, b <-chan T, timeout time.Duration) (v T, from int)
```

## Examples

| Situation | Result | Returns after |
|---|---|---|
| `a` sends `"x"` at 1s, `b` sends `"y"` at 2s, timeout 5s | `"x", 1` | 1s |
| `a` sends at 2s, `b` sends `"y"` at 1s | `"y", 2` | 1s |
| `a` at 3s, `b` at 4s, timeout 2s | `"", 0` | 2s |
| `a` already holds `"x"` in its buffer | `"x", 1` | 0s |
| `a` is `nil`, `b` sends `"y"` at 1s | `"y", 2` | 1s |
| both `nil`, timeout 1s | `"", 0` | 1s |

## Rules
- `from` is `1` for a value from `a`, `2` for a value from `b`, and `0` on timeout (with the zero `T`).
- Receive **one** value only. A `nil` channel never delivers.
- The tests never close `a` or `b`, and never make both deliver at the same instant.
