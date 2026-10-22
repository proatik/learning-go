# Atomic Stats

Count cache hits and misses from many goroutines with `atomic.Int64`: no mutex needed for independent counters.

```go
type Stats struct {
	hits   atomic.Int64
	misses atomic.Int64
}

type Counts struct {
	Hits   int64
	Misses int64
}

func (s *Stats) Hit()
func (s *Stats) Miss()
func (s *Stats) Snapshot() Counts
func (s *Stats) Reset() Counts

func (c Counts) Total() int64
func (c Counts) HitRate() float64
```

## Examples

| Code | Result |
|---|---|
| `var s Stats; s.Snapshot()` | `{Hits:0 Misses:0}` |
| `s.Hit(); s.Hit(); s.Hit(); s.Miss(); s.Snapshot()` | `{Hits:3 Misses:1}` |
| `s.Reset()` | `{Hits:3 Misses:1}`; afterwards `Snapshot()` = `{Hits:0 Misses:0}` |
| `Counts{Hits: 3, Misses: 1}.Total()` | `4` |
| `Counts{Hits: 3, Misses: 1}.HitRate()` | `0.75` |
| `Counts{}.HitRate()` | `0` (no division by zero) |

## Rules
- Use the atomic fields' methods (`Add`, `Load`, `Swap`); no mutex.
- `Reset` returns what each counter held **and** zeroes it in one atomic step per counter. A hit recorded while `Reset` runs must end up either in `Reset`'s result or in the counter, never lost (the tests add up every `Reset` result while 20 goroutines keep counting).
- `Snapshot` and `Reset` read the two counters one after the other, so the pair isn't one atomic moment; that's fine here.
- `Counts` is a plain value type, safe to copy. `Stats` must not be copied (it holds atomics), so its methods use a pointer receiver.
- `HitRate` is `Hits / Total` as a `float64`, and `0` when `Total` is `0`.
