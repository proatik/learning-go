package userapi

import "errors"

// ErrNotFound means the requested user does not exist.
var ErrNotFound = errors.New("not found")

// InputError reports invalid input from the caller.
type InputError struct {
	Field string
	Msg   string
}

// Error returns "invalid <field>: <msg>".
func (e *InputError) Error() string {
	// TODO
	return ""
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
	// TODO
	return User{}, nil
}

// Service holds the business logic on top of a Repo.
type Service struct {
	Repo *Repo
}

// GetUser validates id and loads the user.
func (s *Service) GetUser(id int) (User, error) {
	// TODO
	return User{}, nil
}

// StatusOf maps an error to an HTTP-style status code.
func StatusOf(err error) int {
	// TODO
	return 0
}

// Handle serves a "get user" request.
func Handle(s *Service, id int) (status int, body string) {
	// TODO
	return 0, ""
}
