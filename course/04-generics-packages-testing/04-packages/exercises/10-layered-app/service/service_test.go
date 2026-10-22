package service_test

import (
	"errors"
	"strings"
	"testing"
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/service"
	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/store"
)

var errDB = errors.New("database is down")

// fakeStore is a hand-written test double. Each method fails with errDB
// when its flag is set; otherwise it behaves like a tiny store.
type fakeStore struct {
	tasks                                     []store.Task
	failCreate, failGet, failUpdate, failList bool
	updates                                   int
}

func (f *fakeStore) Create(t store.Task) (store.Task, error) {
	if f.failCreate {
		return store.Task{}, errDB
	}
	t.ID = len(f.tasks) + 1
	f.tasks = append(f.tasks, t)
	return t, nil
}

func (f *fakeStore) Get(id int) (store.Task, error) {
	if f.failGet {
		return store.Task{}, errDB
	}
	if id < 1 || id > len(f.tasks) {
		return store.Task{}, store.ErrNotFound
	}
	return f.tasks[id-1], nil
}

func (f *fakeStore) Update(t store.Task) error {
	if f.failUpdate {
		return errDB
	}
	f.updates++
	f.tasks[t.ID-1] = t
	return nil
}

func (f *fakeStore) List() ([]store.Task, error) {
	if f.failList {
		return nil, errDB
	}
	return append([]store.Task(nil), f.tasks...), nil
}

var _ service.Store = (*fakeStore)(nil)

var t0 = time.Date(2026, 5, 6, 7, 8, 9, 0, time.UTC)

func clock() time.Time { return t0 }

func newService(t *testing.T, f *fakeStore) *service.Service {
	t.Helper()
	s := service.New(f, clock)
	if s == nil {
		t.Fatalf("service.New = nil, want a service")
	}
	return s
}

func TestAddValidation(t *testing.T) {
	long := strings.Repeat("é", service.MaxTitle+1)
	exact := strings.Repeat("é", service.MaxTitle)
	tests := []struct {
		title string
		want  error
	}{
		{"", service.ErrEmptyTitle},
		{"   \t ", service.ErrEmptyTitle},
		{long, service.ErrTitleTooLong},
		{"  " + exact + "  ", nil},
	}
	for _, tt := range tests {
		f := &fakeStore{}
		_, err := newService(t, f).Add(tt.title)
		if !errors.Is(err, tt.want) {
			t.Errorf("Add(%q) error = %v, want %v", tt.title, err, tt.want)
		}
		if tt.want != nil && len(f.tasks) != 0 {
			t.Errorf("Add(%q) stored %d tasks, want 0", tt.title, len(f.tasks))
		}
	}
}

func TestAdd(t *testing.T) {
	f := &fakeStore{}
	s := newService(t, f)
	got, err := s.Add("  Write tests ")
	want := store.Task{ID: 1, Title: "Write tests", Created: t0}
	if err != nil || got != want {
		t.Fatalf("Add(%q) = %+v, %v, want %+v, <nil>", "  Write tests ", got, err, want)
	}
	if _, err := s.Add("write TESTS"); !errors.Is(err, service.ErrDuplicate) {
		t.Errorf("Add(%q) error = %v, want %v", "write TESTS", err, service.ErrDuplicate)
	}
	if len(f.tasks) != 1 {
		t.Fatalf("stored tasks = %d, want 1", len(f.tasks))
	}
	f.tasks[0].Done = true // a done task no longer blocks its title
	if got, err := s.Add("Write tests"); err != nil || got.ID != 2 {
		t.Errorf("Add(%q) after completing = %+v, %v, want ID 2, <nil>", "Write tests", got, err)
	}
}

func TestDefaultClock(t *testing.T) {
	s := service.New(&fakeStore{}, nil)
	if s == nil {
		t.Fatalf("service.New(store, nil) = nil, want a service")
	}
	before := time.Now()
	got, err := s.Add("x")
	if err != nil || got.Created.Before(before) || got.Created.After(time.Now()) {
		t.Errorf("Add(%q) with nil clock Created = %v, %v, want about time.Now()", "x", got.Created, err)
	}
}

func TestComplete(t *testing.T) {
	f := &fakeStore{}
	s := newService(t, f)
	s.Add("a")
	s.Add("b")
	if err := s.Complete(2); err != nil {
		t.Fatalf("Complete(2) error = %v, want <nil>", err)
	}
	if len(f.tasks) != 2 || !f.tasks[1].Done || f.tasks[0].Done {
		t.Fatalf("tasks after Complete(2) = %+v, want only task 2 done", f.tasks)
	}
	if err := s.Complete(2); !errors.Is(err, service.ErrAlreadyDone) {
		t.Errorf("second Complete(2) error = %v, want %v", err, service.ErrAlreadyDone)
	}
	if f.updates != 1 {
		t.Errorf("store updates = %d, want 1", f.updates)
	}
	if err := s.Complete(42); !errors.Is(err, store.ErrNotFound) {
		t.Errorf("Complete(42) error = %v, want %v", err, store.ErrNotFound)
	}
}

func TestPendingAndStats(t *testing.T) {
	f := &fakeStore{}
	s := newService(t, f)
	if got, err := s.Pending(); err != nil || len(got) != 0 {
		t.Errorf("empty Pending() = %+v, %v, want [], <nil>", got, err)
	}
	for _, title := range []string{"a", "b", "c", "d"} {
		s.Add(title)
	}
	s.Complete(1)
	s.Complete(3)
	got, err := s.Pending()
	if err != nil || len(got) != 2 || got[0].Title != "b" || got[1].Title != "d" {
		t.Errorf("Pending() = %+v, %v, want tasks b and d, <nil>", got, err)
	}
	done, pending, err := s.Stats()
	if done != 2 || pending != 2 || err != nil {
		t.Errorf("Stats() = %d, %d, %v, want 2, 2, <nil>", done, pending, err)
	}
}

func TestStoreErrorsAreWrapped(t *testing.T) {
	tests := []struct {
		name string
		f    *fakeStore
		call func(s *service.Service) error
	}{
		{"Add/List fails", &fakeStore{failList: true}, func(s *service.Service) error { _, err := s.Add("x"); return err }},
		{"Add/Create fails", &fakeStore{failCreate: true}, func(s *service.Service) error { _, err := s.Add("x"); return err }},
		{"Complete/Get fails", &fakeStore{failGet: true, tasks: []store.Task{{ID: 1}}}, func(s *service.Service) error { return s.Complete(1) }},
		{"Complete/Update fails", &fakeStore{failUpdate: true, tasks: []store.Task{{ID: 1}}}, func(s *service.Service) error { return s.Complete(1) }},
		{"Pending fails", &fakeStore{failList: true}, func(s *service.Service) error { _, err := s.Pending(); return err }},
		{"Stats fails", &fakeStore{failList: true}, func(s *service.Service) error { _, _, err := s.Stats(); return err }},
	}
	for _, tt := range tests {
		err := tt.call(newService(t, tt.f))
		if !errors.Is(err, errDB) {
			t.Errorf("%s: error = %v, want %v", tt.name, err, errDB)
			continue
		}
		if err == errDB {
			t.Errorf("%s: error = %v (unwrapped), want it wrapped with context", tt.name, err)
		}
	}
}
