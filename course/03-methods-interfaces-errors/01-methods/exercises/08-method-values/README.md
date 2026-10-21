# Method Values and Expressions

Treat methods as ordinary functions: **method values** (`c.Add`, bound to a receiver) and **method expressions** (`(*Calc).Add`, receiver as first argument).

```go
type Calc struct {
	acc int
}

func (c *Calc) Add(n int)
func (c *Calc) Sub(n int)
func (c *Calc) Mul(n int)
func (c *Calc) Result() int

// Bound returns the method values c.Add, c.Sub, c.Mul (in that order).
func Bound(c *Calc) []func(int)

// Lookup returns the method expression for "add", "sub" or "mul", or nil.
func Lookup(name string) func(*Calc, int)

type Step struct {
	Op  string
	Arg int
}

// Run applies the steps to a fresh Calc and returns the result.
func Run(steps []Step) (int, bool)

// RunAll applies the same steps to each calculator in cs.
func RunAll(cs []*Calc, steps []Step) bool
```

## Examples

| Code | Result |
|---|---|
| `var c Calc; ops := Bound(&c); ops[0](5); ops[2](3)` | `c.Result()` = `15` |
| `Lookup("sub")(&c, 4)` | same as `c.Sub(4)` |
| `Lookup("div")` | `nil` |
| `Run([]Step{{"add", 2}, {"mul", 10}, {"sub", 1}})` | `19, true` |
| `Run(nil)` | `0, true` |
| `Run([]Step{{"add", 1}, {"pow", 2}})` | `0, false` |

## Rules
- A new `Calc` starts at `0`; `Add`/`Sub`/`Mul` update the accumulator in place.
- `Bound` must return method values bound to the **given** `c`: calling them changes `c`.
- `Lookup` is case-sensitive: only the exact lower-case names.
- `Run` and `RunAll` must use `Lookup` (method expressions), not a `switch` on each step. An unknown op makes `Run` return `0, false` and `RunAll` return `false` **without changing any** calculator.
