# Parser with Precise Errors

Parse `key=value;key2=value2` lines. On bad input, return a `*ParseError` that points to the **exact line and column**, the way a compiler does.

```go
var (
	ErrMissingEquals = errors.New("missing '='")
	ErrEmptyKey      = errors.New("empty key")
	ErrInvalidKey    = errors.New("invalid key character")
	ErrDuplicateKey  = errors.New("duplicate key")
)

type ParseError struct {
	Line int   // 1-based
	Col  int   // 1-based, counted in runes
	Err  error // one of the sentinels above
}

func (e *ParseError) Error() string // "line 2, col 5: missing '='"
func (e *ParseError) Unwrap() error // e.Err

func Parse(input string) (map[string]string, error)
```

## Examples

| Input | Result |
|---|---|
| `"a=1;b=2"` | `{"a": "1", "b": "2"}` |
| `" host = go.dev ; port=443 ;"` | `{"host": "go.dev", "port": "443"}` |
| `"a=1\n\nb=x=y"` | `{"a": "1", "b": "x=y"}` |
| `"empty="` | `{"empty": ""}` |
| `""` | `{}` (empty, non-nil) |
| `"a=1;oops"` | `line 1, col 5: missing '='` |
| `"a=1\n  =2"` | `line 2, col 3: empty key` |
| `"a=1\nnö=2"` | `line 2, col 2: invalid key character` |
| `"a=1\nb=2;a=3"` | `line 2, col 5: duplicate key` |
| `"name=Zoë;x"` | `line 1, col 10: missing '='` (`ë` is one column) |

## Rules
- Split the input into lines on `\n`, and each line into **segments** on `;`.
- A segment that is empty or only spaces/tabs is skipped. This allows blank lines, `;;` and a trailing `;`.
- In each segment, the **first** `=` separates key from value. Trim spaces and tabs around both. The value may be empty and may contain `=` or any Unicode.
- Keys may only contain ASCII letters, digits, `_`, `.` and `-`.
- Stop at the **first** problem and return `nil, &ParseError{...}`:

| Problem | `Err` | `Col` points at |
|---|---|---|
| no `=` in the segment | `ErrMissingEquals` | the segment's first non-space character |
| key empty after trimming | `ErrEmptyKey` | the `=` |
| bad character in the key (spaces inside a key count) | `ErrInvalidKey` | the first bad character |
| key already defined (on any earlier line) | `ErrDuplicateKey` | the key's first character |

- Columns count **runes**, not bytes, starting at 1 for each line. A tab is one column.
- `Error()` is `fmt.Sprintf("line %d, col %d: %v", Line, Col, Err)`. With `Unwrap`, `errors.Is(err, ErrDuplicateKey)` works.
