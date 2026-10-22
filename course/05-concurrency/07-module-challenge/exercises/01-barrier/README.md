# Reusable Barrier

Build a **barrier**: a meeting point where `n` goroutines wait for each other. Nobody passes until all `n` have arrived; then all of them are released together and the barrier is ready for the next round.

```go
type Barrier struct {
	mu      sync.Mutex
	n       int           // goroutines per round
	count   int           // arrived so far in this round
	round   int           // number of completed rounds
	release chan struct{} // closed to release the current round
}

func New(n int) *Barrier

func (b *Barrier) Wait() int
```

## Examples

| Code (`b := New(3)`) | Result |
|---|---|
| goroutines 1 and 2 call `b.Wait()` | both block |
| goroutine 3 calls `b.Wait()` | all three return `0` (round 0) |
| the same 3 goroutines call `b.Wait()` again | all three block until the third arrives, then return `1` |
| `New(1).Wait()`, three times | never blocks: returns `0`, `1`, `2` |
| `New(0)` | panics: `barrier: n must be at least 1` |

## Rules
- `New` panics with exactly `"barrier: n must be at least 1"` if `n < 1`.
- `Wait` blocks until `n` goroutines (counting the caller) have called `Wait` in the current round, then every one of them returns.
- `Wait` returns the **round number** it belonged to: `0` for the first `n` callers, `1` for the next `n`, and so on.
- The barrier is **reusable**: a goroutine that leaves round `r` and calls `Wait` again joins round `r+1`, even if other goroutines from round `r` have not returned yet.
- Safe for concurrent use; no goroutines of its own.
