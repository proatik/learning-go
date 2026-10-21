# Error Categories with `Is`

Give an error type a custom `Is` method, so one `*HTTPError` can match whole **categories** like "any 4xx" through `errors.Is`.

```go
var (
	ErrClient        = errors.New("client error")         // any 4xx
	ErrServer        = errors.New("server error")         // any 5xx
	ErrRetryable     = errors.New("retryable")            // 429, 502, 503, 504
	ErrInvalidStatus = errors.New("invalid status code")  // outside 100..599
)

type HTTPError struct {
	URL  string
	Code int
}

func (e *HTTPError) Error() string           // "<url>: status <code>"
func (e *HTTPError) Is(target error) bool

func Check(url string, code int) error
func Advice(err error) string
```

## Examples

| Expression | Result |
|---|---|
| `(&HTTPError{"/users/7", 404}).Error()` | `"/users/7: status 404"` |
| `errors.Is(&HTTPError{"/a", 404}, ErrClient)` | `true` |
| `errors.Is(&HTTPError{"/a", 404}, ErrServer)` | `false` |
| `errors.Is(fmt.Errorf("load: %w", &HTTPError{"/a", 503}), ErrServer)` | `true` |
| `errors.Is(&HTTPError{"/a", 503}, ErrRetryable)` | `true` |
| `errors.Is(&HTTPError{"/a", 404}, &HTTPError{"/other", 404})` | `true` (same code) |
| `Check("/a", 204)` | `nil` |
| `Check("/a", 404)` | `&HTTPError{"/a", 404}` |
| `Check("/a", 42)` | `/a: status 42: invalid status code` |
| `Advice(Check("/a", 429))` | `"retry later"` |

## Rules
- `Is(target)` returns `true` when:
  - `target == ErrClient` and `400 <= Code <= 499`
  - `target == ErrServer` and `500 <= Code <= 599`
  - `target == ErrRetryable` and `Code` is `429`, `502`, `503` or `504`
  - `target` is a `*HTTPError` with the same `Code` (the URL is ignored)
  - otherwise `false`. Don't unwrap inside `Is`: `errors.Is` walks the chain itself.
- `Check(url, code)`:
  - `100..399` → `nil`.
  - `400..599` → `&HTTPError{URL: url, Code: code}`.
  - anything else → `fmt.Errorf("%s: status %d: %w", url, code, ErrInvalidStatus)`.
- `Advice(err)`, checked in this order:

| Condition | Result |
|---|---|
| `err == nil` | `ok` |
| matches `ErrRetryable` | `retry later` |
| matches `ErrClient` | `fix the request` |
| matches `ErrServer` | `server problem` |
| otherwise | `unknown error` |
