package todo_test

import (
	"errors"
	"testing"
	"time"

	todo "learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/service"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/store"
)

// The concrete store satisfies the interface declared by its consumer.
var _ service.Store = (*store.Memory)(nil)

func TestEndToEnd(t *testing.T) {
	t0 := time.Date(2026, 3, 4, 5, 6, 7, 0, time.UTC)
	s := todo.New(func() time.Time { return t0 })
	if s == nil {
		t.Fatalf("todo.New = nil, want a service")
	}
	got, err := s.Add("  Write tests ")
	if want := (store.Task{ID: 1, Title: "Write tests", Created: t0}); err != nil || got != want {
		t.Fatalf("Add(%q) = %+v, %v, want %+v, <nil>", "  Write tests ", got, err, want)
	}
	if _, err := s.Add("write TESTS"); !errors.Is(err, service.ErrDuplicate) {
		t.Errorf("Add(%q) error = %v, want %v", "write TESTS", err, service.ErrDuplicate)
	}
	if _, err := s.Add("Ship it"); err != nil {
		t.Fatalf("Add(%q) error = %v, want <nil>", "Ship it", err)
	}
	if err := s.Complete(1); err != nil {
		t.Errorf("Complete(1) error = %v, want <nil>", err)
	}
	if err := s.Complete(1); !errors.Is(err, service.ErrAlreadyDone) {
		t.Errorf("second Complete(1) error = %v, want %v", err, service.ErrAlreadyDone)
	}
	if err := s.Complete(42); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Complete(42) error = %v, want %v", err, store.ErrNotFound)
	}
	if got, err := s.Add("Write tests"); err != nil || got.ID != 3 {
		t.Errorf("Add(%q) after completing = %+v, %v, want ID 3, <nil>", "Write tests", got, err)
	}
	pending, err := s.Pending()
	if err != nil || len(pending) != 2 || pending[0].ID != 2 || pending[1].ID != 3 {
		t.Errorf("Pending() = %+v, %v, want tasks 2 and 3, <nil>", pending, err)
	}
	done, open, err := s.Stats()
	if done != 1 || open != 2 || err != nil {
		t.Errorf("Stats() = %d, %d, %v, want 1, 2, <nil>", done, open, err)
	}
}

func TestAppsAreIndependent(t *testing.T) {
	a, b := todo.New(nil), todo.New(nil)
	if a == nil || b == nil {
		t.Fatalf("todo.New(nil) = nil, want a service")
	}
	if _, err := a.Add("only in a"); err != nil {
		t.Fatalf("Add(%q) error = %v, want <nil>", "only in a", err)
	}
	if _, pending, _ := b.Stats(); pending != 0 {
		t.Errorf("second app Stats() pending = %d, want 0", pending)
	}
}
