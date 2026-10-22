package board

import (
	"cmp"
	"slices"
	"strings"
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
	if every <= 0 {
		every = time.Second
	}
	b := &Board{
		n:      max(n, 1),
		scores: map[string]int{},
		stop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	b.Publish()
	go b.publisher(every)
	return b
}

// publisher calls Publish on every tick until stop is closed.
func (b *Board) publisher(every time.Duration) {
	defer close(b.done)
	tick := time.NewTicker(every)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			b.Publish()
		case <-b.stop:
			return
		}
	}
}

// Add adds points to player's score.
func (b *Board) Add(player string, points int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.scores[player] += points
	b.version++
}

// Transfer moves points from one player to another as a single update.
func (b *Board) Transfer(from, to string, points int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.scores[from] -= points
	b.scores[to] += points
	b.version++
}

// Publish takes a snapshot now and makes it the latest.
func (b *Board) Publish() {
	b.mu.Lock()
	entries := make([]Entry, 0, len(b.scores))
	total := 0
	for p, s := range b.scores {
		entries = append(entries, Entry{p, s})
		total += s
	}
	snap := &Snapshot{Total: total, Players: len(entries), Version: b.version}
	b.mu.Unlock()

	slices.SortFunc(entries, func(x, y Entry) int {
		return cmp.Or(cmp.Compare(y.Score, x.Score), strings.Compare(x.Player, y.Player))
	})
	snap.Top = entries[:min(b.n, len(entries))]
	b.latest.Store(snap)
}

// Latest returns the latest published snapshot without waiting for writers.
func (b *Board) Latest() Snapshot {
	return *b.latest.Load()
}

// Close stops the publishing goroutine and waits for it to exit. It is safe
// to call more than once.
func (b *Board) Close() {
	b.closeOnce.Do(func() { close(b.stop) })
	<-b.done
}
