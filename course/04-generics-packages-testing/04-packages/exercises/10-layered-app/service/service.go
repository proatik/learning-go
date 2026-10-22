// Package service holds the to-do business rules. It depends on the Store
// interface it declares, not on a concrete store.
package service

import (
	"errors"
	"time"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/10-layered-app/store"
)

// MaxTitle is the maximum title length in runes.
const MaxTitle = 80

// Errors returned by Service methods.
var (
	ErrEmptyTitle   = errors.New("service: title is empty")
	ErrTitleTooLong = errors.New("service: title is too long")
	ErrDuplicate    = errors.New("service: an open task has this title")
	ErrAlreadyDone  = errors.New("service: task is already done")
)

// Store is the persistence the service needs. *store.Memory satisfies it,
// and so does any fake a test provides.
type Store interface {
	// TODO: list the four methods Service uses (see the README)
}

// Service applies the to-do rules on top of a Store.
type Service struct {
	store Store
	now   func() time.Time
}

// New returns a Service backed by st. A nil now means time.Now.
func New(st Store, now func() time.Time) *Service {
	// TODO
	return &Service{}
}

// Add creates a task from the trimmed title, after checking that it is
// non-empty, at most MaxTitle runes, and not the title of an open task.
func (s *Service) Add(title string) (store.Task, error) {
	_ = store.Task{} // TODO
	return store.Task{}, nil
}

// Complete marks the task done, or returns ErrAlreadyDone.
func (s *Service) Complete(id int) error {
	// TODO
	return nil
}

// Pending returns the open tasks in ID order.
func (s *Service) Pending() ([]store.Task, error) {
	// TODO
	return nil, nil
}

// Stats counts done and pending tasks.
func (s *Service) Stats() (done, pending int, err error) {
	// TODO
	return 0, 0, nil
}
