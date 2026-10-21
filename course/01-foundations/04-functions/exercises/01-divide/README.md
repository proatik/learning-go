# Divide with Remainder

Implement integer division that reports division by zero as an error instead of crashing:

```go
func Divide(a, b int) (q, r int, err error)
```

| Call | `q` | `r` | `err` |
|---|---|---|---|
| `Divide(17, 5)` | `3` | `2` | `nil` |
| `Divide(-7, 2)` | `-3` | `-1` | `nil` |
| `Divide(0, 4)` | `0` | `0` | `nil` |
| `Divide(5, 0)` | `0` | `0` | non-nil |

## Rules
- Use Go's `/` and `%` (they truncate toward zero, so `-7 / 2 == -3` and `-7 % 2 == -1`).
- On error, `q` and `r` must be `0`.
- Create the error with `errors.New` or `fmt.Errorf`.
- Keep the named results; a naked `return` is fine in a function this short.
