# Request IDs

Store a request ID in a context the safe way, behind an **unexported key type**, and let background work keep the ID after the request is canceled.

```go
type ctxKey struct{}

func WithRequestID(ctx context.Context, id string) context.Context
func RequestID(ctx context.Context) (string, bool)
func Detach(ctx context.Context) context.Context
```

## Examples

| Code | Result |
|---|---|
| `RequestID(WithRequestID(ctx, "req-42"))` | `"req-42", true` |
| `RequestID(context.Background())` | `"", false` |
| `RequestID(WithRequestID(ctx, ""))` | `"", true` (stored, just empty) |
| `RequestID(WithRequestID(WithRequestID(ctx, "outer"), "inner"))` | `"inner", true` |
| `c, cancel := context.WithCancel(WithRequestID(ctx, "a"))`; `RequestID(c)` | `"a", true` (children inherit values) |
| `RequestID(context.WithValue(ctx, "requestID", "evil"))` | `"", false` (a string key is a different key) |
| `d := Detach(req)`, then `req` is canceled | `d.Err()` = `nil`, `d.Done()` = `nil`, `RequestID(d)` still works |

## Rules
- Use `ctxKey{}` as the key: no other package can create a value of that type, so nothing can collide with or overwrite your ID.
- `RequestID` must never panic, whatever else is stored in the context.
- `Detach` returns a context with the same values as `ctx` that is **never** canceled and has **no** deadline (`context.WithoutCancel`).
