# Higher-Order Functions

A *higher-order* function takes functions as arguments or returns one. Using the provided type

```go
type IntFunc func(int) int
```

implement:

```go
func ApplyN(f IntFunc, n, x int) int // f applied n times: f(f(...f(x)))
func Compose(f, g IntFunc) IntFunc   // returns a function computing f(g(x))
func Twice(f IntFunc) IntFunc        // returns a function computing f(f(x))
```

With `inc := func(x int) int { return x + 1 }` and `dbl := func(x int) int { return x * 2 }`:

| Expression | Result |
|---|---|
| `ApplyN(dbl, 3, 1)` | `8` |
| `ApplyN(inc, 0, 42)` | `42` |
| `Compose(inc, dbl)(5)` | `11` (double first, then add one) |
| `Compose(dbl, inc)(5)` | `12` |
| `Twice(inc)(5)` | `7` |
| `Twice(Twice(dbl))(1)` | `16` |

## Rules
- `ApplyN` with `n <= 0` returns `x` unchanged.
- `Compose` and `Twice` return **function literals**; they must not call `f` or `g` themselves until the returned function is called.
- Bonus: write `Twice` in one line using `Compose`.
