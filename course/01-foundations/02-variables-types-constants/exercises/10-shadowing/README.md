# Fix the Shadowing Bugs

`config.go` **compiles** and looks reasonable, but both functions ignore the value they parse and always return the default. The cause is the same in both: a `:=` inside an `if` block declares a **new** variable that shadows the outer one.

```go
func ParsePort(s string) (int, error)            // "" → DefaultPort (8080)
func Scale(value int, factor string) (int, error) // "" → factor 1
```

| Call | Result |
|---|---|
| `ParsePort("")` | `8080, nil` |
| `ParsePort("443")` | `443, nil` |
| `ParsePort("0")` | `0, <error>` (must be 1 … 65535) |
| `ParsePort("http")` | `0, <error>` |
| `Scale(5, "")` | `5, nil` |
| `Scale(5, "3")` | `15, nil` |
| `Scale(5, "-2")` | `-10, nil` |
| `Scale(5, "0")` | `0, <error>` |

## Rules
- Keep the structure. The fix is a few characters per function plus declaring `var err error` where needed.
- Don't rename the variables to dodge the problem; understand which variable each `:=` creates.
- `fmt.Errorf("...")` builds an error value; `return 0, fmt.Errorf(...)` reports a failure.

## Tip
Before fixing, add `fmt.Println(port)` right after the inner `if` block and run the tests with `-v`. Which `port` does it print?
