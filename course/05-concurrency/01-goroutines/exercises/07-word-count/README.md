# Word Count

Count words across many documents, one goroutine per document. A Go map must never be written by two goroutines at once, so each goroutine fills **its own** map and you merge them after `Wait`.

```go
// words splits a document into words. It's given; use it for every document.
// (The tests replace it with a slow version to check documents are counted concurrently.)
var words = strings.Fields

func Count(docs []string) map[string]int
```

## Examples

| `docs` | Result |
|---|---|
| `{"the cat", "The dog the end"}` | `map[cat:1 dog:1 end:1 the:3]` |
| `{"Go go GO"}` | `map[go:3]` |
| `{"  Ünïcode\tünïcode\n"}` | `map[ünïcode:2]` |
| `{"", "   "}` | `map[]` |
| `nil` | `map[]` |

## Rules
- Split each document with `words(doc)` and lowercase every word with `strings.ToLower`.
- Start one goroutine per document. Goroutine `i` counts into a **local** map and stores it in slot `i` of a pre-sized `[]map[string]int`.
- After `Wait`, merge the per-document maps into the result. No goroutine may write the result map.
- Always return a non-nil map (empty when there are no words).
- A shared map written by several goroutines usually crashes the program with `fatal error: concurrent map writes`.
