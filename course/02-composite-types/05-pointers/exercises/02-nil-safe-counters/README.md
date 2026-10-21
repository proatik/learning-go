# Nil-Safe Counters

Write small helpers that work **through** an `*int` and never crash on `nil`.

```go
func Increment(p *int)
func Reset(p *int)
func SafeDeref(p *int, def int) int
func Sum(ps []*int) int
```

| Function | Non-nil `p` | `nil` |
|---|---|---|
| `Increment(p)` | `*p` grows by 1 | does nothing |
| `Reset(p)` | `*p` becomes `0` | does nothing |
| `SafeDeref(p, def)` | returns `*p` | returns `def` |
| `Sum(ps)` | adds every `*p` | nil entries count as `0` |

| Example | Result |
|---|---|
| `n := 5; Increment(&n)` | `n == 6` |
| `SafeDeref(new(0), 9)` | `0` (the pointer is set, its value just happens to be 0) |
| `SafeDeref(nil, 9)` | `9` |
| `Sum([]*int{new(1), nil, new(3)})` | `4` |
| `Sum(nil)` | `0` |

## Rules
- Check `p == nil` **before** any `*p`: dereferencing nil panics.
- `Sum` should reuse `SafeDeref`.
