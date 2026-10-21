package userapi

import (
	"errors"
	"fmt"
	"testing"
)

func newService(fail error) *Service {
	return &Service{Repo: &Repo{
		Users: map[int]User{7: {ID: 7, Name: "Ann"}, 42: {ID: 42, Name: "Zoë"}},
		Fail:  fail,
	}}
}

func TestInputError(t *testing.T) {
	e := &InputError{Field: "id", Msg: "must be positive"}
	if got, want := e.Error(), "invalid id: must be positive"; got != want {
		t.Errorf("(%+v).Error() = %q, want %q", *e, got, want)
	}
}

func TestFind(t *testing.T) {
	r := newService(nil).Repo
	u, err := r.Find(7)
	if err != nil || u != (User{7, "Ann"}) {
		t.Errorf("Find(7) = %+v, %v, want {ID:7 Name:Ann}, <nil>", u, err)
	}

	u, err = r.Find(9)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Find(9) error = %v, want %v", err, ErrNotFound)
	}
	if got, want := err.Error(), "find user 9: not found"; got != want {
		t.Errorf("Find(9) error = %q, want %q", got, want)
	}
	if u != (User{}) {
		t.Errorf("Find(9) = %+v, want the zero User with an error", u)
	}

	dbDown := errors.New("db down")
	r.Fail = dbDown
	_, err = r.Find(7)
	if !errors.Is(err, dbDown) {
		t.Fatalf("Find(7) with Fail error = %v, want %v", err, dbDown)
	}
	if got, want := err.Error(), "find user 7: db down"; got != want {
		t.Errorf("Find(7) with Fail error = %q, want %q", got, want)
	}
}

func TestGetUser(t *testing.T) {
	s := newService(nil)
	u, err := s.GetUser(42)
	if err != nil || u.Name != "Zoë" {
		t.Errorf("GetUser(42) = %+v, %v, want {ID:42 Name:Zoë}, <nil>", u, err)
	}

	_, err = s.GetUser(9)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("GetUser(9) error = %v, want %v", err, ErrNotFound)
	}
	if got, want := err.Error(), "get user: find user 9: not found"; got != want {
		t.Errorf("GetUser(9) error = %q, want %q", got, want)
	}

	for _, id := range []int{0, -3} {
		// Fail is set: a correct GetUser rejects the id before reaching the repo.
		s := newService(errors.New("repo must not be called"))
		u, err := s.GetUser(id)
		ie, ok := errors.AsType[*InputError](err)
		if !ok {
			t.Errorf("GetUser(%d) error = %v, want *InputError", id, err)
			continue
		}
		if *ie != (InputError{Field: "id", Msg: "must be positive"}) {
			t.Errorf("GetUser(%d) error = %+v, want {Field:id Msg:must be positive}", id, *ie)
		}
		if u != (User{}) {
			t.Errorf("GetUser(%d) = %+v, want the zero User with an error", id, u)
		}
	}
}

func TestStatusOf(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{nil, 200},
		{ErrNotFound, 404},
		{fmt.Errorf("a: %w", fmt.Errorf("b: %w", ErrNotFound)), 404},
		{&InputError{Field: "x", Msg: "y"}, 400},
		{fmt.Errorf("ctx: %w", &InputError{Field: "x", Msg: "y"}), 400},
		{errors.New("not found"), 500}, // same text, different value
		{fmt.Errorf("ctx: %v", ErrNotFound), 500},
		{errors.New("boom"), 500},
	}
	for _, tt := range tests {
		if got := StatusOf(tt.err); got != tt.want {
			t.Errorf("StatusOf(%v) = %d, want %d", tt.err, got, tt.want)
		}
	}
}

func TestHandle(t *testing.T) {
	tests := []struct {
		fail       error
		id         int
		wantStatus int
		wantBody   string
	}{
		{nil, 7, 200, "Ann"},
		{nil, 42, 200, "Zoë"},
		{nil, 9, 404, "user not found"},
		{nil, 0, 400, "invalid id: must be positive"},
		{nil, -3, 400, "invalid id: must be positive"},
		{errors.New("db down: password=hunter2"), 7, 500, "internal error"},
		{fmt.Errorf("timeout: %w", ErrNotFound), 7, 404, "user not found"},
	}
	for _, tt := range tests {
		status, body := Handle(newService(tt.fail), tt.id)
		if status != tt.wantStatus || body != tt.wantBody {
			t.Errorf("Handle(%d) with Fail=%v = %d, %q, want %d, %q", tt.id, tt.fail, status, body, tt.wantStatus, tt.wantBody)
		}
	}
}
