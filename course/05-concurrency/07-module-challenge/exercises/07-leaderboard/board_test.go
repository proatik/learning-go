package board

import (
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func wantSnap(t *testing.T, step string, got Snapshot, want Snapshot) {
	t.Helper()
	if got.Top == nil {
		t.Fatalf("%s: Latest().Top = nil, want %v (empty, not nil)", step, want.Top)
	}
	if !slices.Equal(got.Top, want.Top) || got.Total != want.Total || got.Players != want.Players || got.Version != want.Version {
		t.Fatalf("%s: Latest() = %+v, want %+v", step, got, want)
	}
}

func TestReadmeScenario(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(2, time.Second)
		defer b.Close()
		wantSnap(t, "New(2, 1s)", b.Latest(), Snapshot{Top: []Entry{}})
		b.Add("ann", 5)
		wantSnap(t, `Add("ann", 5), before any publish`, b.Latest(), Snapshot{Top: []Entry{}})
		time.Sleep(time.Second)
		synctest.Wait()
		wantSnap(t, "1s later (ticker published)", b.Latest(),
			Snapshot{Top: []Entry{{"ann", 5}}, Total: 5, Players: 1, Version: 1})
		b.Add("bob", 7)
		b.Add("cat", 5)
		b.Publish()
		wantSnap(t, `Add("bob", 7), Add("cat", 5), Publish()`, b.Latest(),
			Snapshot{Top: []Entry{{"bob", 7}, {"ann", 5}}, Total: 17, Players: 3, Version: 3})
		b.Transfer("bob", "cat", 3)
		b.Publish()
		wantSnap(t, `Transfer("bob", "cat", 3), Publish()`, b.Latest(),
			Snapshot{Top: []Entry{{"cat", 8}, {"ann", 5}}, Total: 17, Players: 3, Version: 4})
		b.Add("dan", 5)
		b.Add("cat", -3)
		b.Publish()
		wantSnap(t, `Add("dan", 5), Add("cat", -3), Publish()`, b.Latest(),
			Snapshot{Top: []Entry{{"ann", 5}, {"cat", 5}}, Total: 19, Players: 4, Version: 6})
	})
}

func TestTopSizes(t *testing.T) {
	tests := []struct {
		n    int
		want []Entry
	}{
		{1, []Entry{{"eve", 9}}},
		{3, []Entry{{"eve", 9}, {"bob", 4}, {"dan", 4}}},
		{10, []Entry{{"eve", 9}, {"bob", 4}, {"dan", 4}, {"ann", 1}, {"cat", -2}}},
		{0, []Entry{{"eve", 9}}},
		{-5, []Entry{{"eve", 9}}},
	}
	for _, tt := range tests {
		synctest.Test(t, func(t *testing.T) {
			b := New(tt.n, time.Hour)
			defer b.Close()
			b.Add("dan", 4)
			b.Add("ann", 1)
			b.Add("eve", 9)
			b.Add("cat", -2)
			b.Transfer("eve", "bob", 0) // creates bob at 0
			b.Add("bob", 4)
			b.Publish()
			wantSnap(t, fmt.Sprintf("New(%d), 5 players", tt.n), b.Latest(),
				Snapshot{Top: tt.want, Total: 16, Players: 5, Version: 6})
		})
	}
}

func TestPeriodicPublish(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(1, time.Minute)
		defer b.Close()
		b.Add("x", 1)
		time.Sleep(59 * time.Second)
		synctest.Wait()
		if v := b.Latest().Version; v != 0 {
			t.Fatalf("New(1, 1m): Latest().Version after 59s = %d, want 0", v)
		}
		time.Sleep(time.Second)
		synctest.Wait()
		if v := b.Latest().Version; v != 1 {
			t.Fatalf("New(1, 1m): Latest().Version after 1m = %d, want 1", v)
		}
		b.Add("x", 1)
		time.Sleep(time.Minute)
		synctest.Wait()
		if v := b.Latest().Version; v != 2 {
			t.Fatalf("New(1, 1m): Latest().Version after 2m = %d, want 2", v)
		}
	})
}

func TestDefaultInterval(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(1, 0)
		defer b.Close()
		b.Add("x", 1)
		time.Sleep(time.Second)
		synctest.Wait()
		if v := b.Latest().Version; v != 1 {
			t.Fatalf("New(1, 0): Latest().Version after 1s = %d, want 1 (default every 1s)", v)
		}
	})
}

func TestCloseStopsPublisher(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := New(1, time.Second)
		b.Close()
		b.Close() // a second Close must not panic or block
		b.Add("x", 1)
		time.Sleep(5 * time.Second)
		synctest.Wait()
		if v := b.Latest().Version; v != 0 {
			t.Fatalf("Latest().Version 5s after Close = %d, want 0 (no more ticks)", v)
		}
		b.Publish()
		if v := b.Latest().Version; v != 1 {
			t.Fatalf("Publish() after Close: Latest().Version = %d, want 1", v)
		}
	})
}

func TestLatestDoesNotLock(t *testing.T) {
	b := New(1, time.Hour)
	defer b.Close()
	b.mu.Lock()
	defer b.mu.Unlock()
	done := make(chan struct{})
	go func() {
		b.Latest()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatalf("Latest() blocked while a writer held mu, want it to return at once")
	}
}

func TestSnapshotsAreImmutable(t *testing.T) {
	b := New(3, time.Hour)
	defer b.Close()
	b.Add("a", 3)
	b.Add("b", 2)
	b.Add("c", 1)
	b.Publish()
	old := b.Latest()
	keep := slices.Clone(old.Top)
	b.Add("c", 10)
	b.Add("d", 20)
	b.Publish()
	if !slices.Equal(old.Top, keep) {
		t.Errorf("an old snapshot's Top changed to %v after a new Publish, want %v", old.Top, keep)
	}
	if want := []Entry{{"d", 20}, {"c", 11}, {"a", 3}}; !slices.Equal(b.Latest().Top, want) {
		t.Errorf("Latest().Top = %v, want %v", b.Latest().Top, want)
	}
}

func TestNoHalfAppliedUpdates(t *testing.T) {
	const players, writers, transfers = 20, 8, 2000
	b := New(5, time.Millisecond)
	defer b.Close()
	for i := range players {
		b.Add(fmt.Sprintf("p%02d", i), 100)
	}
	var stop atomic.Bool
	var mu sync.Mutex
	var problem string
	report := func(s string) {
		mu.Lock()
		defer mu.Unlock()
		if problem == "" {
			problem = s
		}
	}
	finished := make(chan struct{})
	go func() {
		var writersWG, readers sync.WaitGroup
		for w := range writers {
			writersWG.Go(func() {
				for i := range transfers {
					from := fmt.Sprintf("p%02d", (w*7+i)%players)
					to := fmt.Sprintf("p%02d", (w*3+i*11+1)%players)
					b.Transfer(from, to, 1+i%5)
				}
			})
		}
		for range 4 {
			readers.Go(func() {
				last := 0
				for !stop.Load() {
					b.Publish()
					s := b.Latest()
					switch {
					case s.Total != players*100:
						report(fmt.Sprintf("Latest().Total = %d, want %d (a Transfer was half-applied)", s.Total, players*100))
						return
					case s.Players != players:
						report(fmt.Sprintf("Latest().Players = %d, want %d", s.Players, players))
						return
					case s.Version < last:
						report(fmt.Sprintf("Latest().Version went back from %d to %d", last, s.Version))
						return
					case len(s.Top) != 5:
						report(fmt.Sprintf("len(Latest().Top) = %d, want 5", len(s.Top)))
						return
					}
					last = s.Version
				}
			})
		}
		writersWG.Wait()
		stop.Store(true)
		readers.Wait()
		close(finished)
	}()
	select {
	case <-finished:
	case <-time.After(2 * time.Second):
		t.Fatalf("%d writers x %d Transfers with 4 readers did not finish within 2s", writers, transfers)
	}
	if problem != "" {
		t.Fatal(problem)
	}
	b.Publish()
	if v, want := b.Latest().Version, players+writers*transfers; v != want {
		t.Errorf("Latest().Version after all updates = %d, want %d (lost updates?)", v, want)
	}
}
