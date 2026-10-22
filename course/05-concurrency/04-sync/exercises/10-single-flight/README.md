# Single Flight

When 100 requests for `user:42` arrive at once and the cache is empty, don't hit the database 100 times. A `Group` lets the **first** caller run `fn` while everyone else asking for the same key waits and gets the same result. Services use this to protect databases from a "thundering herd".

```go
type call[V any] struct {
	done chan struct{} // closed when fn has returned and val, err are set
	val  V
	err  error
	dups int // how many callers joined this call instead of starting their own
}

type Group[V any] struct {
	mu    sync.Mutex          // guards calls
	calls map[string]*call[V] // in-flight calls by key
}

func (g *Group[V]) Do(key string, fn func() (V, error)) (v V, shared bool, err error)
```

## Examples

| Code | Result |
|---|---|
| `var g Group[int]; g.Do("k", fn)` (nobody else asking) | `fn`'s value, `false`, `fn`'s error |
| the same again, after it returned | calls `fn` **again** (results aren't cached) |
| 10 goroutines `g.Do("k", fn)` while `fn` runs | `fn` runs **once**; all 10 get its value and error, with `shared` = `true` |
| `g.Do("a", …)` and `g.Do("b", …)` at the same time | both `fn`s run in parallel |
| a failing `fn` | its error goes to every caller that shared the call |

## Rules
- The zero `Group` is ready to use: create `g.calls` on first use.
- **First caller** for a key: store a new `*call[V]` (with a `done` channel) in `g.calls`, **unlock**, run `fn`, then remove the key from `g.calls` and close `done`.
- **Later callers** while that call is in flight: increment its `dups` (under `g.mu`), unlock, wait on `<-c.done`, return `c.val`, `true`, `c.err`.
- `shared` is `true` for every caller of a call with `dups > 0`, the first caller included; `false` when one caller had it alone.
- Never hold `g.mu` while `fn` runs, or while waiting on `done`: other keys must not wait, and joiners must be able to register. The tests inspect `g.calls` and `dups` (using `g.mu.TryLock`) while `fn` is blocked.
- After `Do` returns, the key is gone from `g.calls`: the next `Do` for it calls `fn` again.
- You may assume `fn` doesn't panic.
- The real thing lives in `golang.org/x/sync/singleflight` (same idea; it returns `(v, err, shared)`).
