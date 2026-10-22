package bank

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// run calls f in a goroutine and fails the test if it doesn't return within
// 2s, so a self-deadlock fails the test instead of hanging it.
func run(t *testing.T, what string, f func()) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		f()
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not finish within 2s (deadlock?)", what)
	}
}

func TestSequential(t *testing.T) {
	var a Account
	steps := []struct {
		op      string
		amount  int
		wantErr error
		want    int
	}{
		{"Deposit", 100, nil, 100},
		{"Withdraw", 30, nil, 70},
		{"Withdraw", 70, nil, 0},
		{"Withdraw", 1, ErrInsufficientFunds, 0},
		{"Deposit", 0, ErrInvalidAmount, 0},
		{"Deposit", -5, ErrInvalidAmount, 0},
		{"Deposit", 50, nil, 50},
		{"Withdraw", 51, ErrInsufficientFunds, 50},
		{"Withdraw", 0, ErrInvalidAmount, 50},
		{"Withdraw", -10, ErrInvalidAmount, 50},
		{"Withdraw", 50, nil, 0},
	}
	for _, s := range steps {
		var err error
		var got int
		run(t, s.op, func() {
			if s.op == "Deposit" {
				err = a.Deposit(s.amount)
			} else {
				err = a.Withdraw(s.amount)
			}
			got = a.Balance()
		})
		switch {
		case s.wantErr == nil && err != nil:
			t.Errorf("%s(%d) error = %v, want nil", s.op, s.amount, err)
		case s.wantErr != nil && !errors.Is(err, s.wantErr):
			t.Errorf("%s(%d) error = %v, want %v", s.op, s.amount, err, s.wantErr)
		}
		if got != s.want {
			t.Errorf("%s(%d): Balance() = %d, want %d", s.op, s.amount, got, s.want)
		}
	}
}

func TestConcurrentDeposits(t *testing.T) {
	var a Account
	run(t, "50 goroutines × 1000 Deposit(2)", func() {
		var wg sync.WaitGroup
		for range 50 {
			wg.Go(func() {
				for range 1000 {
					a.Deposit(2)
				}
			})
		}
		wg.Wait()
	})
	if got := a.Balance(); got != 100000 {
		t.Errorf("50 goroutines × 1000 Deposit(2): Balance() = %d, want 100000", got)
	}
}

func TestNoOverdraft(t *testing.T) {
	for round := range 20 {
		var a Account
		a.balance = 500
		var ok, insufficient atomic.Int64
		run(t, "1000 concurrent Withdraw(1)", func() {
			var wg sync.WaitGroup
			for range 1000 {
				wg.Go(func() {
					switch err := a.Withdraw(1); {
					case err == nil:
						ok.Add(1)
					case errors.Is(err, ErrInsufficientFunds):
						insufficient.Add(1)
					}
				})
			}
			wg.Wait()
		})
		if got := ok.Load(); got != 500 {
			t.Fatalf("round %d: balance 500, 1000 concurrent Withdraw(1): successes = %d, want 500", round, got)
		}
		if got := insufficient.Load(); got != 500 {
			t.Fatalf("round %d: balance 500, 1000 concurrent Withdraw(1): ErrInsufficientFunds count = %d, want 500", round, got)
		}
		if got := a.Balance(); got != 0 {
			t.Fatalf("round %d: balance 500, 1000 concurrent Withdraw(1): Balance() = %d, want 0", round, got)
		}
	}
}

func TestMixedConcurrent(t *testing.T) {
	var a Account
	a.balance = 1000
	var withdrawn atomic.Int64
	run(t, "concurrent deposits and withdrawals", func() {
		var wg sync.WaitGroup
		for g := range 40 {
			wg.Go(func() {
				for range 500 {
					if g%2 == 0 {
						a.Deposit(3)
					} else if a.Withdraw(5) == nil {
						withdrawn.Add(5)
					}
					if b := a.Balance(); b < 0 {
						t.Errorf("Balance() = %d during the run, want >= 0", b)
						return
					}
				}
			})
		}
		wg.Wait()
	})
	want := 1000 + 20*500*3 - int(withdrawn.Load())
	if got := a.Balance(); got != want {
		t.Errorf("start 1000 + deposits %d - successful withdrawals %d: Balance() = %d, want %d", 20*500*3, withdrawn.Load(), got, want)
	}
}
