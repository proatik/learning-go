# Retry Temporary Failures

Retry an operation only when its error says it is **temporary**, a behaviour checked through an interface, not through a specific type.

```go
var ErrNoAttempts = errors.New("attempts must be at least 1")

// TemporaryError marks any error as temporary.
type TemporaryError struct {
	Err error
}

func (e *TemporaryError) Error() string   // e.Err.Error()
func (e *TemporaryError) Temporary() bool // true
func (e *TemporaryError) Unwrap() error   // e.Err

func IsTemporary(err error) bool
func Retry(attempts int, f func() error) error
```

## Examples

| `f` returns, call by call | `Retry(3, f)` returns | calls |
|---|---|---|
| `nil` | `nil` | 1 |
| temp, temp, `nil` | `nil` | 3 |
| temp `"busy"`, temp `"busy"`, temp `"busy"` | `gave up after 3 attempts: busy` | 3 |
| temp, `errors.New("bad key")` | `attempt 2: bad key` | 2 |
| `errors.New("bad key")` | `attempt 1: bad key` | 1 |

`Retry(0, f)` returns `ErrNoAttempts` and never calls `f`.

## Rules
- `IsTemporary(err)`:
  - Find the **first** error in `err`'s tree that has a `Temporary() bool` method, and return what that method says.
  - `false` if there is none (including `err == nil`).
  - Use `errors.As` with a pointer to an interface variable, or `errors.AsType` with an interface type argument.
- `Retry(attempts, f)`:
  - `attempts < 1` → return `ErrNoAttempts` itself, without calling `f`.
  - Call `f` up to `attempts` times. Return `nil` as soon as it succeeds.
  - If attempt `i` fails with a **non-temporary** error, stop right away: `fmt.Errorf("attempt %d: %w", i, err)`.
  - If every attempt failed temporarily: `fmt.Errorf("gave up after %d attempts: %w", attempts, lastErr)`.
  - Always wrap with `%w`, so callers can still use `errors.Is` / `errors.As` on the cause.
- `TemporaryError` methods are one-liners. `Error` returns the inner message unchanged.
