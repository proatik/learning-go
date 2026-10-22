# Bounded Downloader

Download a list of URLs the way a real service should: at most `limit` requests in flight, a **timeout per request**, **one retry** for errors that might go away, and an overall `ctx` that can cancel the whole batch. Results come back in input order.

```go
type Client interface {
	Get(ctx context.Context, url string) ([]byte, error)
}

var ErrTemporary = errors.New("fetch: temporary failure")

type Result struct {
	URL      string
	Body     []byte
	Err      error
	Attempts int // Get calls made for this URL: 0, 1 or 2
}

func Download(ctx context.Context, c Client, urls []string, limit int, timeout time.Duration) []Result
```

## Examples

| Situation (`limit` 3, `timeout` 2s) | `Result` |
|---|---|
| `Get` returns `"ok"` | `{URL, Body: "ok", Err: nil, Attempts: 1}` |
| first `Get` fails with an error wrapping `ErrTemporary`, second returns `"ok"` | `{Body: "ok", Err: nil, Attempts: 2}` |
| both attempts fail with `ErrTemporary` | `{Err: <second error>, Attempts: 2}` |
| `Get` fails with `ErrNotFound` (not temporary) | `{Err: ErrNotFound, Attempts: 1}`: no retry |
| `Get` takes 5s | first attempt times out after 2s, retried, times out again: `{Err: context.DeadlineExceeded, Attempts: 2}` after 4s |
| 10 URLs of 1s each | never more than 3 `Get` calls at once; done after 4s |
| `ctx` cancelled while some URLs are in flight and others haven't started | in flight: the error `Get` returned, no retry; not started: `{Err: ctx.Err(), Attempts: 0}` |

## Rules
- Return one `Result` per URL, **in input order** (`results[i].URL == urls[i]`), even for duplicates. Empty input → empty result.
- At most `limit` `Get` calls run at the same time (`limit < 1` counts as `1`). A retry reuses the same slot.
- Every attempt gets its **own** context: `context.WithTimeout(ctx, timeout)` (always `cancel` it).
- **Retry once** if the attempt failed with `ErrTemporary` (`errors.Is`) or with `context.DeadlineExceeded` from its own timeout, but only while `ctx` itself is still alive (`ctx.Err() == nil`). Never retry anything else, and never more than once.
- If `ctx` is done before a URL gets a slot, don't call `Get` for it: `Err = ctx.Err()`, `Attempts = 0`.
- `Download` returns only when every goroutine it started has finished.
