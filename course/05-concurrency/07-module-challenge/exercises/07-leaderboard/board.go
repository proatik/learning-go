package board

import (
	"sync"
	"sync/atomic"
	"time"
)

// Entry is one player's score.
type Entry struct {
	Player string
	Score  int
}

// Snapshot is an immutable view of the board at one moment.
type Snapshot struct {
	Top     []Entry // best first, at most n entries
	Total   int     // sum of every player's score
	Players int     // number of players
	Version int     // updates applied when the snapshot was taken
}

// Board keeps scores and publishes snapshots of them. Create one with New
// and stop its publisher with Close.
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

// New returns an empty board whose snapshots list the top n players, and
// starts a goroutine that publishes a snapshot every every.
func New(n int, every time.Duration) *Board {
	// TODO
	return &Board{}
}

// Add adds points to player's score.
func (b *Board) Add(player string, points int) {
	// TODO
}

// Transfer moves points from one player to another as a single update.
func (b *Board) Transfer(from, to string, points int) {
	// TODO
}

// Publish takes a snapshot now and makes it the latest.
func (b *Board) Publish() {
	// TODO
}

// Latest returns the latest published snapshot without waiting for writers.
func (b *Board) Latest() Snapshot {
	// TODO
	return Snapshot{}
}

// Close stops the publishing goroutine and waits for it to exit. It is safe
// to call more than once.
func (b *Board) Close() {
	// TODO
}
