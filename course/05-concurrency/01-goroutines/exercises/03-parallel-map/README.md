# Parallel Map

Apply `f` to every element **concurrently**, one goroutine per element, and return the results **in input order**.

```go
func Map[T, U any](in []T, f func(T) U) []U
```

## Examples

| Call | Result |
|---|---|
| `Map([]int{3, 1, 2}, f)` where `f(n)` sleeps `n` seconds, returns `n * 10` | `[30 10 20]`, after **3s** (not 6s) |
| `Map([]string{"go", "héllo"}, utf8.RuneCountInString)` | `[2 5]` |
| `Map([]int{}, f)` / `Map(nil, f)` | an empty slice; `f` is never called |

## Rules
- Allocate the result with `make([]U, len(in))` **before** starting any goroutine; goroutine `i` writes only `out[i]`.
- `out[i]` must be `f(in[i])`, whatever order the goroutines finish in.
- Call `f` exactly once per element, and return only after every call has finished.
- Never `append` to a shared slice from a goroutine.
