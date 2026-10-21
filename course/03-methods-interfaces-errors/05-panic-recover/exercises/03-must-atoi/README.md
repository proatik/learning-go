# MustAtoi

Write `Must`-style helpers that **panic** with a clear `error` value instead of returning one, like `regexp.MustCompile`.

```go
var ErrPortRange = errors.New("port out of range")

func Must(n int, err error) int  // n, or panic(err)
func MustAtoi(s string) int      // strconv.Atoi, or panic
func MustPort(s string) int      // a TCP port 1..65535, or panic
```

| Call | Result |
|---|---|
| `Must(42, nil)` | `42` |
| `Must(0, io.EOF)` | panics with `io.EOF` itself |
| `MustAtoi("-17")` | `-17` |
| `MustAtoi("12a")` | panics with error `MustAtoi("12a"): strconv.Atoi: parsing "12a": invalid syntax` |
| `MustPort("8080")` | `8080` |
| `MustPort("70000")` | panics with error `MustPort("70000"): port out of range` |
| `MustPort("http")` | panics with error `MustAtoi("http"): strconv.Atoi: parsing "http": invalid syntax` |

## Rules
- `Must`: if `err != nil`, `panic(err)` with the **same** error value; otherwise return `n`.
- `MustAtoi`: the panic value is an `error` that wraps the `strconv` error with `%w` (so `errors.Is(v, strconv.ErrSyntax)` / `strconv.ErrRange` work).
- `MustPort`: parse with `MustAtoi` (its panic passes through unchanged), then panic with `fmt.Errorf("MustPort(%q): %w", s, ErrPortRange)` unless `1 <= n <= 65535`.
- Successful calls must **not** panic.
