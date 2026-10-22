package dag

import (
	"context"
	"errors"
)

// Job is one unit of work that may depend on other jobs by name.
type Job struct {
	Name string
	Deps []string
	Run  func(ctx context.Context) error
}

// Status is the outcome of one job.
type Status int

const (
	Skipped Status = iota // never ran
	Done                  // ran and returned nil
	Failed                // ran and returned an error
)

// String returns "skipped", "done" or "failed".
func (s Status) String() string {
	switch s {
	case Skipped:
		return "skipped"
	case Done:
		return "done"
	case Failed:
		return "failed"
	}
	return "unknown"
}

// Validation errors returned by Run before any job starts.
var (
	ErrDuplicate  = errors.New("dag: duplicate job")
	ErrUnknownDep = errors.New("dag: unknown dependency")
	ErrCycle      = errors.New("dag: dependency cycle")
)

// Run validates the jobs, then runs each one as soon as all of its
// dependencies are done. It returns every job's status and the first job
// failure.
func Run(ctx context.Context, jobs []Job) (map[string]Status, error) {
	// TODO
	return nil, nil
}
