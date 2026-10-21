# Numbers ↔ Strings with `strconv`

Implement:

```go
func Double(s string) (int, error)            // parse s as an int, return 2×
func AddStrings(a, b string) (float64, error) // parse both as float64, return the sum
func Label(n int) string                      // "item-" followed by n
func Binary(n int) string                     // n in base 2
```

| Call | Result |
|---|---|
| `Double("21")` | `42, nil` |
| `Double("-7")` | `-14, nil` |
| `Double("abc")` | `0, <error>` |
| `Double("1.5")` | `0, <error>` (not an integer) |
| `AddStrings("1.5", "2.25")` | `3.75, nil` |
| `AddStrings("1", "x")` | `0, <error>` |
| `Label(7)` | `"item-7"` |
| `Binary(10)` | `"1010"` |
| `Binary(-5)` | `"-101"` |

## Rules
- Use `strconv.Atoi`, `strconv.ParseFloat(s, 64)`, `strconv.Itoa` and `strconv.FormatInt`.
- On failure return the **zero value** plus the error you got from `strconv` (don't make up your own).
- `string(n)` does **not** produce digits. Don't use it.
