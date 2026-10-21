# Report Every Problem

A signup form should list **all** its problems at once, not one per submit. Collect them and combine them with `errors.Join`.

```go
type Form struct {
	Username string
	Email    string
	Password string
	Confirm  string
	Age      int
	Terms    bool
}

var (
	ErrUsernameRequired = errors.New("username is required")
	ErrUsernameInvalid  = errors.New("username must be 3-16 letters, digits or _")
	ErrEmailInvalid     = errors.New("invalid email")
	ErrPasswordShort    = errors.New("password must be at least 8 characters")
	ErrPasswordNoDigit  = errors.New("password must contain a digit")
	ErrPasswordMismatch = errors.New("passwords do not match")
	ErrUnderage         = errors.New("must be at least 18")
	ErrTerms            = errors.New("terms must be accepted")
)

func Validate(f Form) error
func Messages(err error) []string
```

## Examples

| Form | `Validate` matches (`errors.Is`) |
|---|---|
| `{"ann_99", "ann@go.dev", "s3cretpass", "s3cretpass", 30, true}` | nothing: returns `nil` |
| `{"", "ann@go.dev", "s3cretpass", "s3cretpass", 30, true}` | `ErrUsernameRequired` |
| `{"a!", "nope", "short", "other", 12, false}` | `ErrUsernameInvalid`, `ErrEmailInvalid`, `ErrPasswordShort`, `ErrPasswordNoDigit`, `ErrPasswordMismatch`, `ErrUnderage`, `ErrTerms` |

For `f := Form{"ann", "nope", "abcdefgh1", "abcdefgh1", 17, true}`:

```go
Validate(f).Error()     // "email \"nope\": invalid email\nage 17: must be at least 18"
Messages(Validate(f))   // []string{`email "nope": invalid email`, "age 17: must be at least 18"}
```

## Rules
- Run the checks in this order and keep **every** failure:

| # | Check | Error added |
|---|---|---|
| 1 | `Username == ""` | `ErrUsernameRequired` (skip check 2) |
| 2 | not 3–16 runes, or a char other than ASCII letter, digit, `_` | `ErrUsernameInvalid` |
| 3 | `Email` not exactly one `@` with a non-empty part on each side, or no `.` in the domain part | `fmt.Errorf("email %q: %w", f.Email, ErrEmailInvalid)` |
| 4 | `Password` shorter than 8 runes | `ErrPasswordShort` |
| 5 | `Password` has no digit `0-9` | `ErrPasswordNoDigit` |
| 6 | `Confirm != Password` | `ErrPasswordMismatch` |
| 7 | `Age < 18` | `fmt.Errorf("age %d: %w", f.Age, ErrUnderage)` |
| 8 | `!Terms` | `ErrTerms` |

- No problems: return `nil`. `errors.Join` with no arguments already does this.
- `Messages(err)`:
  - `nil` → `nil`.
  - an error with an `Unwrap() []error` method (what `errors.Join` returns) → the `Error()` of each part, in order.
  - any other error → a one-element slice with its `Error()`.
