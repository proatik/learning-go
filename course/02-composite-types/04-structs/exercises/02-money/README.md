# Money in Cents

Floating-point can't represent `0.10` exactly, so money is stored as an **integer number of cents** wrapped in a struct.

```go
type Money struct {
	Cents int64
}

func New(dollars, cents int64) Money  // dollars*100 + cents
func Add(a, b Money) Money
func Sum(ms []Money) Money
func Format(m Money) string           // "$12.05", "-$0.50"
func Split(m Money, n int) []Money    // divide into n parts that add up to m
```

| Call | Result |
|---|---|
| `New(12, 5)` | `Money{Cents: 1205}` |
| `New(-1, -50)` | `Money{Cents: -150}` |
| `Add(Money{150}, Money{275})` | `Money{Cents: 425}` |
| `Sum(nil)` | `Money{}` |
| `Format(Money{1205})` | `"$12.05"` |
| `Format(Money{7})` | `"$0.07"` |
| `Format(Money{-50})` | `"-$0.50"` |
| `Split(Money{100}, 3)` | `[{34} {33} {33}]` |
| `Split(Money{5}, 0)` | `nil` |

## Rules
- `Format`: a `-` before the `$` for negative amounts, always exactly two cent digits, no thousands separators.
- `Split`: for `n <= 0` return `nil`. Otherwise return exactly `n` parts whose sum is `m`; parts differ by at most one cent and the **first** parts get the extra cents. You may assume `m.Cents >= 0`.
- No `float64` anywhere.
