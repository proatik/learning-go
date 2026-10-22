// Package store persists tasks. Memory is the in-memory implementation;
// a database implementation could be added later with the same methods.
package store

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"time"
)

// ErrNotFound is wrapped by Get and Update when no task has the given ID.
var ErrNotFound = errors.New("store: not found")

// Task is a to-do item.
type Task struct {
	ID      int
	Title   string
	Done    bool
	Created time.Time
}

// Memory stores tasks in a map. Create one with NewMemory.
type Memory struct {
	tasks  map[int]Task
	lastID int
}

// NewMemory returns an empty in-memory store.
func NewMemory() *Memory {
	return &Memory{tasks: map[int]Task{}}
}

// Create stores t under the next free ID (1, 2, 3, …), ignoring t.ID, and
// returns the stored task.
func (m *Memory) Create(t Task) (Task, error) {
	m.lastID++
	t.ID = m.lastID
	m.tasks[t.ID] = t
	return t, nil
}

// Get returns the task with the given ID, or an error wrapping ErrNotFound.
func (m *Memory) Get(id int) (Task, error) {
	t, ok := m.tasks[id]
	if !ok {
		return Task{}, fmt.Errorf("%w: task %d", ErrNotFound, id)
	}
	return t, nil
}

// Update replaces the stored task with the same ID, or returns an error
// wrapping ErrNotFound.
func (m *Memory) Update(t Task) error {
	if _, ok := m.tasks[t.ID]; !ok {
		return fmt.Errorf("%w: task %d", ErrNotFound, t.ID)
	}
	m.tasks[t.ID] = t
	return nil
}

// List returns all tasks sorted by ID.
func (m *Memory) List() ([]Task, error) {
	tasks := make([]Task, 0, len(m.tasks))
	for _, id := range slices.Sorted(maps.Keys(m.tasks)) {
		tasks = append(tasks, m.tasks[id])
	}
	return tasks, nil
}
