package main

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"sync"
	"time"
)

const (
	StatusNotStarted = ""
	StatusAttempted  = "attempted"
	StatusPassed     = "passed"
)

type ExerciseProgress struct {
	Status        string     `json:"status,omitempty"`
	HintsShown    int        `json:"hintsShown,omitempty"`
	SolutionShown bool       `json:"solutionShown,omitempty"`
	GaveUp        bool       `json:"gaveUp,omitempty"` // solution viewed before passing
	Runs          int        `json:"runs,omitempty"`
	PassedAt      *time.Time `json:"passedAt,omitempty"`
}

// ProgressStore persists learner progress to a JSON file.
type ProgressStore struct {
	mu    sync.Mutex
	path  string
	items map[string]*ExerciseProgress
}

func openProgress(path string) (*ProgressStore, error) {
	s := &ProgressStore{path: path, items: map[string]*ExerciseProgress{}}
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &s.items); err != nil {
		return nil, err
	}
	return s, nil
}

// Get returns a copy of the progress for an exercise.
func (s *ProgressStore) Get(id string) ExerciseProgress {
	s.mu.Lock()
	defer s.mu.Unlock()
	if p, ok := s.items[id]; ok {
		return *p
	}
	return ExerciseProgress{}
}

func (s *ProgressStore) Update(id string, fn func(p *ExerciseProgress)) (ExerciseProgress, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.items[id]
	if !ok {
		p = &ExerciseProgress{}
		s.items[id] = p
	}
	fn(p)
	return *p, s.save()
}

func (s *ProgressStore) ResetAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.items)
	return s.save()
}

func (s *ProgressStore) save() error {
	b, err := json.MarshalIndent(s.items, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
