# Paginated API

Hide pagination behind an iterator: the caller ranges over items, and pages are fetched only when needed. Errors travel through an `iter.Seq2[T, error]`.

```go
// Fetch loads one page (numbered from 1). more reports whether a later page exists.
type Fetch[T any] func(page int) (items []T, more bool, err error)

func Items[T any](fetch Fetch[T]) iter.Seq2[T, error]
func CollectAll[T any](seq iter.Seq2[T, error]) ([]T, error)
```

## Examples

A fetch returning page 1 `[a b] more`, page 2 `[] more`, page 3 `[c] done`:

| Code | Result |
|---|---|
| `for item, err := range Items(fetch)` | `(a, nil)`, `(b, nil)`, `(c, nil)`; pages 1, 2, 3 fetched |
| `break` after `a` | Only page 1 is fetched |
| `break` after `b` | Only page 1 is fetched |
| `CollectAll(Items(fetch))` | `[a b c], nil` |

A fetch where page 2 fails with `err`:

| Code | Result |
|---|---|
| `for item, err := range Items(fetch)` | page 1's items with `nil`, then `("", err)` once, then the loop ends |
| `CollectAll(Items(fetch))` | page 1's items, `err` |

## Rules
- `Items` calls `fetch(1)`, `fetch(2)`, … in order. It yields every item of a page with a `nil` error before fetching the next page.
- Fetch the next page **only** when the current page is used up and `more` was `true`. An empty page with `more == true` is allowed: just fetch the next one.
- When `fetch` returns an error, yield `(zero value, err)` once and stop. Ignore that call's `items` and `more`.
- Building the iterator fetches nothing. A `break` stops fetching. Each `range` starts again from page 1.
- `CollectAll` returns every item in order and `nil`, or the items received before the first error and that error. It must stop at the first error.
