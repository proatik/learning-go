# Parse an Age

Parse a user-typed age and report **which** problem occurred using sentinel errors.

```go
var (
	ErrEmpty      = errors.New("empty age")
	ErrNotNumber  = errors.New("not a number")
	ErrOutOfRange = errors.New("out of range")
)

func ParseAge(s string) (int, error)
```

## Examples

| Call | Result | `err.Error()` |
|---|---|---|
| `ParseAge("42")` | `42, nil` | |
| `ParseAge("  7 ")` | `7, nil` | |
| `ParseAge("")` | `0, ErrEmpty` | `empty age` |
| `ParseAge("   ")` | `0, ErrEmpty` | `empty age` |
| `ParseAge("abc")` | `0, err` (`errors.Is(err, ErrNotNumber)`) | `parse age "abc": not a number` |
| `ParseAge("151")` | `0, err` (`errors.Is(err, ErrOutOfRange)`) | `age 151: out of range` |
| `ParseAge("-1")` | `0, err` (`errors.Is(err, ErrOutOfRange)`) | `age -1: out of range` |

## Rules
- First trim the surrounding spaces with `strings.TrimSpace`. Use the trimmed string everywhere below.
- Empty after trimming: return `ErrEmpty` itself.
- Not an integer (`strconv.Atoi` fails): wrap → `parse age "<s>": not a number` (`%q` of the trimmed input).
- Valid ages are `0` to `150` inclusive. Otherwise wrap → `age <n>: out of range`.
- Return `0` with every error.
- Wrap with `%w`: the tests use `errors.Is`.
