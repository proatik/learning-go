# Calculator with Function Values

Pick an operation **as a value**, then call it.

```go
type BinOp func(a, b int) int

func Lookup(op string) (BinOp, error)
func Calc(op string, a, b int) (int, error)
```

`Lookup` maps an operator to a function:

| `op` | Function returned |
|---|---|
| `"+"` | `a + b` |
| `"-"` | `a - b` |
| `"*"` | `a * b` |
| `"/"` | `a / b` |
| `"%"` | `a % b` |
| anything else | `nil` and a non-nil error |

`Calc` looks up the operation and applies it:

| Call | Result |
|---|---|
| `Calc("+", 2, 3)` | `5, nil` |
| `Calc("-", 2, 3)` | `-1, nil` |
| `Calc("/", 7, 2)` | `3, nil` |
| `Calc("/", 7, 0)` | `0, <error>` |
| `Calc("%", 7, 0)` | `0, <error>` |
| `Calc("^", 2, 3)` | `0, <error>` |

## Rules
- `Lookup` uses a `switch` and returns named functions or function literals.
- `Calc` must call `Lookup` and pass its error on unchanged (`return 0, err`).
- Division or remainder by zero is an error, checked in `Calc` **before** calling the function (a real `7 / 0` on ints crashes the program).
