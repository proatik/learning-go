# Number Bases

Implement three functions in `bases.go`:

```go
func ToBase(n, base int) string
func FromBase(s string, base int) (int, error)
func Describe(n int) string
```

| Call | Result |
|---|---|
| `ToBase(10, 2)` | `"1010"` |
| `ToBase(255, 16)` | `"ff"` |
| `ToBase(-8, 8)` | `"-10"` |
| `FromBase("ff", 16)` | `255, nil` |
| `FromBase("1010", 2)` | `10, nil` |
| `FromBase("12", 2)` | `0, error` |
| `Describe(255)` | `"255 = 0b11111111 = 0o377 = 0xFF"` |
| `Describe(0)` | `"0 = 0b0 = 0o0 = 0x0"` |

## Rules
- `ToBase` must use `strconv.FormatInt` (bases 2–36, lowercase digits).
- `FromBase` must use `strconv.ParseInt`. On failure return `0` and the error from `strconv`.
- `Describe` must use **one** `fmt.Sprintf` call with the verbs `%d`, `%b`, `%o` and `%X`. Bonus: pass `n` only once by using `%[1]d`.
