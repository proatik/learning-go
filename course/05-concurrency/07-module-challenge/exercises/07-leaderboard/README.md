# Live Leaderboard

A game server updates player scores from many goroutines while many more goroutines read the leaderboard. Writers share the scores under a mutex; readers get an **immutable snapshot**, swapped in atomically, so a read never waits for a writer and never sees an update half-applied.

```go
type Entry struct {
	Player string
	Score  int
}

type Snapshot struct {
	Top     []Entry // best first, at most n entries
	Total   int     // sum of every player's score
	Players int     // number of players
	Version int     // updates applied when the snapshot was taken
}

type Board struct {
	mu        sync.Mutex
	n         int
	scores    map[string]int
	version   int
	latest    atomic.Pointer[Snapshot]
	stop      chan struct{}
	done      chan struct{}
	closeOnce sync.Once
}

func New(n int, every time.Duration) *Board

func (b *Board) Add(player string, points int)
func (b *Board) Transfer(from, to string, points int)
func (b *Board) Publish()
func (b *Board) Latest() Snapshot
func (b *Board) Close()
```

## Examples

| Code (`b := New(2, time.Second)`) | `b.Latest()` |
|---|---|
| right after `New` | `{Top: [], Total: 0, Players: 0, Version: 0}` |
| `b.Add("ann", 5)`, at once | unchanged: nothing published yet |
| 1s later (the ticker published) | `{Top: [{ann 5}], Total: 5, Players: 1, Version: 1}` |
| `b.Add("bob", 7)`; `b.Add("cat", 5)`; `b.Publish()` | `{Top: [{bob 7} {ann 5}], Total: 17, Players: 3, Version: 3}` |
| `b.Transfer("bob", "cat", 3)`; `b.Publish()` | `{Top: [{cat 8} {ann 5}], Total: 17, Players: 3, Version: 4}` |
| `b.Add("dan", 5)`; `b.Add("cat", -3)`; `b.Publish()` | `{Top: [{ann 5} {cat 5}], Total: 19, Players: 4, Version: 6}` (ties by name) |

## Rules
- `New` sets `n` (`< 1` counts as `1`), publishes an empty snapshot (`Top` empty, not nil), and starts a goroutine that calls `Publish` every `every` (a `time.Ticker`; `every <= 0` counts as one second).
- `Add` adds `points` (may be negative) to `player`, creating the player at `0` if new. `Transfer` moves `points` from `from` to `to` (creating either at `0` if new). Each call is **one** update: `Version` goes up by one, and no snapshot can ever see only half of a `Transfer`.
- `Publish` takes a snapshot **now** and makes it the latest. `Top` holds the `n` best players: higher score first, ties by name (A to Z). Do the sorting **outside** the lock.
- `Latest` returns the latest published snapshot **without locking `mu`** (load the `atomic.Pointer`). A published snapshot is never modified afterwards.
- `Close` stops the publishing goroutine and waits until it has exited; calling it again does nothing. `Add`, `Transfer`, `Publish` and `Latest` keep working after `Close`.
- Safe for concurrent use.
