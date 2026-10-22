package store_test

import (
	"errors"
	"testing"
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/store"
)

var t0 = time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

func TestCreateAndGet(t *testing.T) {
	m := store.NewMemory()
	for i, title := range []string{"a", "b", "c"} {
		got, err := m.Create(store.Task{ID: 99, Title: title, Created: t0})
		want := store.Task{ID: i + 1, Title: title, Created: t0}
		if err != nil || got != want {
			t.Fatalf("Create(%q) = %+v, %v, want %+v, <nil>", title, got, err, want)
		}
	}
	got, err := m.Get(2)
	if want := (store.Task{ID: 2, Title: "b", Created: t0}); err != nil || got != want {
		t.Errorf("Get(2) = %+v, %v, want %+v, <nil>", got, err, want)
	}
}

func TestNotFound(t *testing.T) {
	m := store.NewMemory()
	for _, id := range []int{0, 1, -1} {
		got, err := m.Get(id)
		if !errors.Is(err, store.ErrNotFound) || got != (store.Task{}) {
			t.Errorf("Get(%d) = %+v, %v, want {}, %v", id, got, err, store.ErrNotFound)
		}
	}
	if err := m.Update(store.Task{ID: 5, Title: "x"}); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Update(ID 5) error = %v, want %v", err, store.ErrNotFound)
	}
	_, err := m.Get(7)
	if want := "store: not found: task 7"; err == nil || err.Error() != want {
		t.Errorf("Get(7) error = %v, want %q", err, want)
	}
}

func TestUpdateAndList(t *testing.T) {
	m := store.NewMemory()
	if got, err := m.List(); err != nil || len(got) != 0 {
		t.Errorf("empty List() = %+v, %v, want [], <nil>", got, err)
	}
	for _, title := range []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven"} {
		if _, err := m.Create(store.Task{Title: title}); err != nil {
			t.Fatalf("Create(%q) error = %v, want <nil>", title, err)
		}
	}
	if err := m.Update(store.Task{ID: 2, Title: "TWO", Done: true}); err != nil {
		t.Fatalf("Update(ID 2) error = %v, want <nil>", err)
	}
	all, err := m.List()
	if err != nil || len(all) != 11 {
		t.Fatalf("List() = %d tasks, %v, want 11, <nil>", len(all), err)
	}
	for i, task := range all {
		if task.ID != i+1 {
			t.Errorf("List()[%d].ID = %d, want %d (sorted by ID)", i, task.ID, i+1)
		}
	}
	if got := all[1]; got.Title != "TWO" || !got.Done {
		t.Errorf("List()[1] = %+v, want Title %q, Done true", got, "TWO")
	}
}
