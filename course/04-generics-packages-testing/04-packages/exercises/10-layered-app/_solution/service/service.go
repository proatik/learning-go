// Package service holds the to-do business rules. It depends on the Store
// interface it declares, not on a concrete store.
package service

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

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
	Create(t store.Task) (store.Task, error)
	Get(id int) (store.Task, error)
	Update(t store.Task) error
	List() ([]store.Task, error)
}

// Service applies the to-do rules on top of a Store.
type Service struct {
	store Store
	now   func() time.Time
}

// New returns a Service backed by st. A nil now means time.Now.
func New(st Store, now func() time.Time) *Service {
	if now == nil {
		now = time.Now
	}
	return &Service{store: st, now: now}
}

// Add creates a task from the trimmed title, after checking that it is
// non-empty, at most MaxTitle runes, and not the title of an open task.
func (s *Service) Add(title string) (store.Task, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return store.Task{}, ErrEmptyTitle
	}
	if utf8.RuneCountInString(title) > MaxTitle {
		return store.Task{}, ErrTitleTooLong
	}
	open, err := s.Pending()
	if err != nil {
		return store.Task{}, fmt.Errorf("add: %w", err)
	}
	for _, t := range open {
		if strings.EqualFold(t.Title, title) {
			return store.Task{}, ErrDuplicate
		}
	}
	t, err := s.store.Create(store.Task{Title: title, Created: s.now()})
	if err != nil {
		return store.Task{}, fmt.Errorf("add: %w", err)
	}
	return t, nil
}

// Complete marks the task done, or returns ErrAlreadyDone.
func (s *Service) Complete(id int) error {
	t, err := s.store.Get(id)
	if err != nil {
		return fmt.Errorf("complete %d: %w", id, err)
	}
	if t.Done {
		return ErrAlreadyDone
	}
	t.Done = true
	if err := s.store.Update(t); err != nil {
		return fmt.Errorf("complete %d: %w", id, err)
	}
	return nil
}

// Pending returns the open tasks in ID order.
func (s *Service) Pending() ([]store.Task, error) {
	all, err := s.store.List()
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}
	var open []store.Task
	for _, t := range all {
		if !t.Done {
			open = append(open, t)
		}
	}
	return open, nil
}

// Stats counts done and pending tasks.
func (s *Service) Stats() (done, pending int, err error) {
	all, err := s.store.List()
	if err != nil {
		return 0, 0, fmt.Errorf("stats: %w", err)
	}
	for _, t := range all {
		if t.Done {
			done++
		} else {
			pending++
		}
	}
	return done, pending, nil
}
