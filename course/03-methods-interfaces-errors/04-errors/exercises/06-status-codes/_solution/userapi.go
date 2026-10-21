package userapi

import (
	"errors"
	"fmt"
)

// ErrNotFound means the requested user does not exist.
var ErrNotFound = errors.New("not found")

// InputError reports invalid input from the caller.
type InputError struct {
	Field string
	Msg   string
}

// Error returns "invalid <field>: <msg>".
func (e *InputError) Error() string {
	return "invalid " + e.Field + ": " + e.Msg
}

// User is a stored user.
type User struct {
	ID   int
	Name string
}

// Repo is an in-memory user store.
type Repo struct {
	Users map[int]User
	Fail  error // if non-nil, every Find fails with it
}

// Find returns the user with the given id.
func (r *Repo) Find(id int) (User, error) {
	if r.Fail != nil {
		return User{}, fmt.Errorf("find user %d: %w", id, r.Fail)
	}
	u, ok := r.Users[id]
	if !ok {
		return User{}, fmt.Errorf("find user %d: %w", id, ErrNotFound)
	}
	return u, nil
}

// Service holds the business logic on top of a Repo.
type Service struct {
	Repo *Repo
}

// GetUser validates id and loads the user.
func (s *Service) GetUser(id int) (User, error) {
	if id <= 0 {
		return User{}, &InputError{Field: "id", Msg: "must be positive"}
	}
	u, err := s.Repo.Find(id)
	if err != nil {
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return u, nil
}

// StatusOf maps an error to an HTTP-style status code.
func StatusOf(err error) int {
	if err == nil {
		return 200
	}
	if errors.Is(err, ErrNotFound) {
		return 404
	}
	if _, ok := errors.AsType[*InputError](err); ok {
		return 400
	}
	return 500
}

// Handle serves a "get user" request.
func Handle(s *Service, id int) (status int, body string) {
	u, err := s.GetUser(id)
	status = StatusOf(err)
	switch status {
	case 200:
		return status, u.Name
	case 404:
		return status, "user not found"
	case 400:
		ie, _ := errors.AsType[*InputError](err)
		return status, ie.Error()
	default:
		return status, "internal error"
	}
}
