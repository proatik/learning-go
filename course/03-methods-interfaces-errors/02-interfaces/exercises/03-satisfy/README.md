# Satisfy the Contract

The interfaces are fixed. The starter types *almost* satisfy them: wrong names, wrong receivers, wrong result types, missing methods. Fix the types, then lock the result in with compile-time assertions.

```go
type Namer interface{ Name() string }
type Labeler interface{ Label() string }
type Counter interface {
	Inc()
	Count() int
}
type Resetter interface{ Reset() }
type ResettableCounter interface {
	Counter
	Resetter
}

type Tag struct{ Text string }  // must satisfy Namer and Labeler as a VALUE
type Tally struct{ n int }      // *Tally must satisfy ResettableCounter
type Celsius float64            // must satisfy Labeler as a value
```

## Examples

| Call | Result |
|---|---|
| `Tag{"go"}.Name()` | `"go"` |
| `Tag{"go"}.Label()` | `"#go"` |
| `t := &Tally{}; t.Inc(); t.Inc(); t.Count()` | `2` |
| `t.Reset(); t.Count()` | `0` |
| `Celsius(21.5).Label()` | `"21.5°C"` |
| `Celsius(-40).Label()` | `"-40.0°C"` |

## Rules
- Don't change the interfaces or the type definitions; change only the **methods**.
- `Tag` (not just `*Tag`) must satisfy `Namer` and `Labeler`.
- `Inc` must really change the tally: a `Tally` used through the interface must count.
- `Celsius.Label` formats with exactly one decimal followed by `°C`.
- Add these lines at the end of `contracts.go` once everything compiles with them:
  ```go
  var (
  	_ Namer             = Tag{}
  	_ Labeler           = Tag{}
  	_ ResettableCounter = (*Tally)(nil)
  	_ Labeler           = Celsius(0)
  )
  ```
- The tests check each contract at run time with type assertions, so a broken type fails its test instead of the whole build.
