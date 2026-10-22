# MapReduce

Build a tiny in-process **MapReduce**: `m` mapper goroutines turn documents into key/value pairs, each pair is routed to one of `r` reducer goroutines by a hash of its key, and each reducer combines the values of its keys. Then use it to count words.

```go
func Partition(key string, r int) int

func MapReduce(
	docs []string,
	m, r int,
	mapFn func(doc string, emit func(key string, value int)),
	reduceFn func(key string, values []int) int,
) map[string]int

func WordCount(docs []string, m, r int) map[string]int
```

## Examples

| Call | Result |
|---|---|
| `Partition("go", 4)` | `3` (FNV-1a hash of `"go"` mod 4) |
| `Partition("anything", 1)` | `0` |
| `WordCount([]string{"the cat", "The dog the end"}, 2, 3)` | `map[cat:1 dog:1 end:1 the:3]` |
| `WordCount(nil, 4, 4)` | `map[]` (empty, not nil) |
| `MapReduce(docs, 3, 2, emitLengths, maxOf)` | for each key, `maxOf` of every value emitted for it |

## Rules
- **`Partition`**: hash `key` with 32-bit FNV-1a (`hash/fnv`'s `New32a`) and return `int(sum % uint32(r))`. `r < 1` counts as `1`.
- **Mappers**: start `m` goroutines (`m < 1` counts as `1`) that take documents from a shared channel and call `mapFn(doc, emit)`. `mapFn` is called exactly once per document.
- **Shuffle**: `emit(key, value)` sends the pair to reducer `Partition(key, r)` over that reducer's own channel. Close the reducer channels only after **every** mapper has finished.
- **Reducers**: start `r` goroutines (`r < 1` counts as `1`). Each one collects the values of its keys into its **own** `map[string][]int` (no locks needed: only that goroutine touches it). When its channel is closed, it calls `reduceFn(key, values)` **once per key** and keeps the results.
- Values reach `reduceFn` in no particular order.
- **Result**: merge every reducer's results into one map (each key lives in exactly one reducer, so nothing collides). Never nil. Every goroutine has finished when `MapReduce` returns.
- **`WordCount`**: a word is an element of `strings.Fields(strings.ToLower(doc))`; map emits `(word, 1)`, reduce sums.
