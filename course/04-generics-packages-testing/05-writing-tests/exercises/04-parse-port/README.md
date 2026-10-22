# Error paths

Write `TestParsePort` that checks the happy path **and every error path** of `ParsePort`, using `errors.Is`.

```go
var (
	ErrEmpty      = errors.New("empty port")
	ErrNotNumber  = errors.New("port is not a number")
	ErrOutOfRange = errors.New("port out of range")
)

// ParsePort parses a TCP port from 1 to 65535.
func ParsePort(s string) (int, error)
```

`ParsePort` is already written (in the read-only checker). You write the test in `port_test.go`.

## Examples

| Call | Port | Error (`errors.Is`) |
|---|---|---|
| `ParsePort("8080")` | `8080` | `nil` |
| `ParsePort("")` | `0` | `ErrEmpty` |
| `ParsePort("http")` | `0` | `ErrNotNumber` |
| `ParsePort("+80")` | `0` | `ErrNotNumber` |
| `ParsePort("0")` | `0` | `ErrOutOfRange` |
| `ParsePort("70000")` | `0` | `ErrOutOfRange` |

## Rules

- Empty input: `ErrEmpty`.
- Any character that is not `0`–`9` (a sign, a space, a letter): `ErrNotNumber`.
- A number outside `1`–`65535`: `ErrOutOfRange`.
- **On any error the port is `0`.**
- Errors are wrapped with the input, e.g. `port "http": port is not a number`. Compare with `errors.Is`, never `==`.
- Report mismatches as `"ParsePort(%q) error = %v, want %v"` and `"ParsePort(%q) = %d, want %d"`.
- **The checker** runs your test on the correct `ParsePort` (it must pass), then once per planted bug (it must fail every time).
