# Fetch All

Fetch many keys from a slow service **concurrently**, like a microservice calling another one for a batch of IDs, and build a report of what succeeded and what failed.

```go
type Fetcher interface {
	Fetch(key string) (string, error)
}

var ErrEmptyKey = errors.New("fetch: empty key")

type Report struct {
	Values map[string]string // key -> value, for every successful fetch
	Errors map[string]error  // key -> error, for every failed key
	Failed []string          // the keys of Errors, sorted ascending
}

func FetchAll(f Fetcher, keys []string) Report
```

## Examples

With a fetcher where `"a"` → `"apple"`, `"b"` → `"banana"`, `"x"` fails with `ErrNotFound`, `"p"` panics with `"boom"`, and every call takes 1s:

| `keys` | Report |
|---|---|
| `{"b", "a"}` | `Values` = `{a: apple, b: banana}`, `Errors` = `{}`, `Failed` = `[]`; after **1s** |
| `{"x", "a", "x"}` | `Values` = `{a: apple}`, `Errors` = `{x: fetch "x": not found}`, `Failed` = `[x]`; `Fetch("x")` called **once** |
| `{"p", "", "a"}` | `Values` = `{a: apple}`, `Errors` = `{"": fetch: empty key, p: fetch "p" panicked: boom}`, `Failed` = `["" p]` |
| `nil` | empty (non-nil) maps, empty `Failed` |

## Rules
- Fetch every **distinct** key exactly once (duplicates in `keys` are ignored), each in its own goroutine, all at the same time.
- The empty key `""` is never fetched: its error is `ErrEmptyKey` itself.
- A fetch error `err` becomes `fmt.Errorf("fetch %q: %w", key, err)` (so `errors.Is(…, err)` works).
- A panic in `Fetch` must not crash the program: recover it in that goroutine and use `fmt.Errorf("fetch %q panicked: %v", key, r)`.
- Every distinct key ends up in exactly one of `Values` or `Errors`.
- **Build both maps after `Wait`**, from per-key results stored by index; goroutines must not write to a map.
- `Values` and `Errors` are never nil. `Failed` is sorted ascending (`""` sorts first).
