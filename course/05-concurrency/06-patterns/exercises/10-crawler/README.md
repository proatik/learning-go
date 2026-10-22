# Concurrent Web Crawler

Crawl a link graph concurrently, starting from one URL: fetch each page once, follow its links up to a depth, never run more than `limit` fetches at once, and stop cleanly when the context ends.

```go
type Fetcher interface {
	Fetch(ctx context.Context, url string) (links []string, err error)
}

func Crawl(ctx context.Context, f Fetcher, start string, depth, limit int) ([]string, error)
```

```text
depth 0:  home
depth 1:  a  b            ← fetched concurrently, at most `limit` at once
depth 2:  c  d            ← links already seen (home, b) are skipped
```

## Examples

The test site: `home → a, b`, `a → b, c, a`, `b → home, d`, `c → e`, `d → (none)`, `e → home, c`.

| Code | Result |
|---|---|
| `Crawl(ctx, f, "home", 0, 3)` | `["home"]`, `nil` |
| `Crawl(ctx, f, "home", 1, 3)` | `["a" "b" "home"]`, `nil` |
| `Crawl(ctx, f, "home", 3, 3)` | `["a" "b" "c" "d" "e" "home"]`, `nil` |
| `Crawl(ctx, f, "home", -1, 3)` | `[]`, `nil`, and no fetches |
| `home → a, zz-missing, m-missing` and `a → m-missing` | `["a" "home"]`, error `"fetch m-missing: not found\nfetch zz-missing: not found"` |
| `home` with 20 linked pages of 1s each, `limit 4` | at most 4 fetches at once; done after 6s |
| a chain of 1s pages with a 2.5s timeout | `["a" "home"]` after exactly 2.5s; `errors.Is(err, context.DeadlineExceeded)` |

## Rules
- `start` is at depth 0. Fetch every URL whose **shortest** link distance from `start` is at most `depth`. A negative depth fetches nothing.
- Fetch every URL **at most once**, even when many pages link to it at the same moment. Keep a `seen` set guarded by a `sync.Mutex`, and check-and-mark it in one locked step.
- Crawl **level by level**: start all fetches of one depth concurrently, wait for them, then move on to the newly found links. This guarantees the shortest-path rule even when some pages are slow.
- At most `max(limit, 1)` fetches run at the same time (a semaphore). Acquire before `go`.
- A failed fetch is not visited and its links aren't followed. Wrap each failure as `fmt.Errorf("fetch %s: %w", url, err)`. Return them joined with `errors.Join`, **sorted by URL**.
- Pass `ctx` to every `Fetch`. Once `ctx` is done, start no new fetches, wait for the running ones, and append `ctx.Err()` to the joined errors.
- Return the visited URLs sorted (empty when nothing was fetched).
