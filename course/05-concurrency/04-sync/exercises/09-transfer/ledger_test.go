package ledger

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// start runs f in a goroutine and returns a channel closed when f returns.
func start(f func()) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	return done
}

// finishes fails the test if done isn't closed within 2s.
func finishes(t *testing.T, done <-chan struct{}, what string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

// waits reports whether done is still open after 20ms.
func waits(done <-chan struct{}) bool {
	time.Sleep(20 * time.Millisecond)
	select {
	case <-done:
		return false
	default:
		return true
	}
}

func TestNewAccount(t *testing.T) {
	a := NewAccount(7, 250)
	if a == nil {
		t.Fatalf("NewAccount(7, 250) = nil, want an account")
	}
	if a.ID != 7 || a.balance != 250 {
		t.Errorf("NewAccount(7, 250) = {ID:%d balance:%d}, want {ID:7 balance:250}", a.ID, a.balance)
	}
	var got int
	finishes(t, start(func() { got = a.Balance() }), "Balance()")
	if got != 250 {
		t.Errorf("NewAccount(7, 250).Balance() = %d, want 250", got)
	}
}

func TestTransferSequential(t *testing.T) {
	a := &Account{ID: 1, balance: 100}
	b := &Account{ID: 2, balance: 50}
	steps := []struct {
		name         string
		from, to     *Account
		amount       int
		wantErr      error
		wantA, wantB int
	}{
		{"a→b 30", a, b, 30, nil, 70, 80},
		{"b→a 80 (everything)", b, a, 80, nil, 150, 0},
		{"b→a 1 (empty)", b, a, 1, ErrInsufficientFunds, 150, 0},
		{"a→b 151", a, b, 151, ErrInsufficientFunds, 150, 0},
		{"a→b 0", a, b, 0, ErrInvalidAmount, 150, 0},
		{"a→b -5", a, b, -5, ErrInvalidAmount, 150, 0},
		{"a→a 10", a, a, 10, ErrSameAccount, 150, 0},
		{"a→b 150", a, b, 150, nil, 0, 150},
	}
	for _, s := range steps {
		var err error
		finishes(t, start(func() { err = Transfer(s.from, s.to, s.amount) }), "Transfer "+s.name)
		switch {
		case s.wantErr == nil && err != nil:
			t.Errorf("Transfer %s error = %v, want nil", s.name, err)
		case s.wantErr != nil && !errors.Is(err, s.wantErr):
			t.Errorf("Transfer %s error = %v, want %v", s.name, err, s.wantErr)
		}
		if a.balance != s.wantA || b.balance != s.wantB {
			t.Errorf("Transfer %s: balances = (%d, %d), want (%d, %d)", s.name, a.balance, b.balance, s.wantA, s.wantB)
		}
	}
}

func TestTotalSequential(t *testing.T) {
	a := &Account{ID: 3, balance: 10}
	b := &Account{ID: 1, balance: 20}
	c := &Account{ID: 2, balance: -5}
	tests := []struct {
		name     string
		accounts []*Account
		want     int
	}{
		{"Total()", nil, 0},
		{"Total(a)", []*Account{a}, 10},
		{"Total(a, b, c)", []*Account{a, b, c}, 25},
		{"Total(c, b)", []*Account{c, b}, 15},
	}
	for _, tt := range tests {
		got := -1
		finishes(t, start(func() { got = Total(tt.accounts...) }), tt.name)
		if got != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, got, tt.want)
		}
	}
	if a.ID != 3 || b.ID != 1 || c.ID != 2 {
		t.Errorf("Total must not change the accounts")
	}
}

func TestTransferLocksInIDOrder(t *testing.T) {
	low := &Account{ID: 3, balance: 100}
	high := &Account{ID: 7, balance: 100}
	for _, dir := range []struct {
		name     string
		from, to *Account
	}{
		{"Transfer(high, low, 10)", high, low},
		{"Transfer(low, high, 10)", low, high},
	} {
		low.mu.Lock() // the lock every transfer must take first
		done := start(func() { Transfer(dir.from, dir.to, 10) })
		if !waits(done) {
			low.mu.Unlock()
			t.Errorf("%s while account 3 is locked = returned, want it to wait for the lock", dir.name)
			continue
		}
		if high.mu.TryLock() {
			high.mu.Unlock()
		} else {
			t.Errorf("%s locked account 7 while waiting for account 3, want account 3 (smaller ID) locked first", dir.name)
		}
		low.mu.Unlock()
		finishes(t, done, dir.name+" after unlocking")
	}
	if low.balance != 100 || high.balance != 100 {
		t.Errorf("after one transfer each way: balances = (%d, %d), want (100, 100)", low.balance, high.balance)
	}
}

func TestTotalLocksEverythingInIDOrder(t *testing.T) {
	a1 := &Account{ID: 1, balance: 1}
	a2 := &Account{ID: 2, balance: 2}
	a3 := &Account{ID: 3, balance: 3}
	a1.mu.Lock()
	got := -1
	done := start(func() { got = Total(a3, a1, a2) })
	if !waits(done) {
		a1.mu.Unlock()
		t.Fatalf("Total(a3, a1, a2) while account 1 is locked = returned, want it to wait for the lock")
	}
	for _, a := range []*Account{a2, a3} {
		if a.mu.TryLock() {
			a.mu.Unlock()
		} else {
			t.Errorf("Total(a3, a1, a2) locked account %d while waiting for account 1, want account 1 (smallest ID) locked first", a.ID)
		}
	}
	a1.mu.Unlock()
	finishes(t, done, "Total(a3, a1, a2) after unlocking")
	if got != 6 {
		t.Errorf("Total(a3, a1, a2) = %d, want 6", got)
	}
	for _, a := range []*Account{a1, a2, a3} {
		if !a.mu.TryLock() {
			t.Fatalf("after Total returned, account %d is still locked, want every lock released", a.ID)
		}
		a.mu.Unlock()
	}
}

func TestConcurrentTransfers(t *testing.T) {
	accounts := []*Account{
		{ID: 4, balance: 1000},
		{ID: 1, balance: 1000},
		{ID: 3, balance: 1000},
		{ID: 2, balance: 1000},
	}
	const want = 4000
	var badTotal, negative atomic.Int64
	stop := make(chan struct{})
	var checker sync.WaitGroup
	checker.Go(func() { // the money in the system never changes
		for {
			if Total(accounts[2], accounts[0], accounts[3], accounts[1]) != want {
				badTotal.Add(1)
			}
			select {
			case <-stop:
				return
			default:
			}
		}
	})
	var wg sync.WaitGroup
	for g := range 40 {
		wg.Go(func() {
			for i := range 1000 {
				from := accounts[(g+i)%4]
				to := accounts[(g+i+1+g%3)%4] // includes both directions of every pair
				switch err := Transfer(from, to, 1+i%7); {
				case err == nil, errors.Is(err, ErrInsufficientFunds):
				default:
					t.Errorf("Transfer(%d → %d) error = %v, want nil or %v", from.ID, to.ID, err, ErrInsufficientFunds)
					return
				}
				if from.Balance() < 0 {
					negative.Add(1)
				}
			}
		})
	}
	finishes(t, start(wg.Wait), "40 goroutines × 1000 transfers in both directions")
	close(stop)
	finishes(t, start(checker.Wait), "the Total checker")
	if got := Total(accounts...); got != want {
		t.Errorf("after 40000 concurrent transfers: Total() = %d, want %d", got, want)
	}
	if n := badTotal.Load(); n != 0 {
		t.Errorf("Total() during the transfers was not %d on %d calls, want 0 (lock every account, transfer atomically)", want, n)
	}
	if n := negative.Load(); n != 0 {
		t.Errorf("a balance went below 0 %d times, want 0", n)
	}
}
