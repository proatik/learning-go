package dag

import (
	"context"
	"errors"
	"fmt"
	"sync"
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

// node is a job plus its outcome. status is written before done is closed,
// so anyone who has received from done may read it.
type node struct {
	job    Job
	done   chan struct{}
	status Status
}

// Run validates the jobs, then runs each one as soon as all of its
// dependencies are done. It returns every job's status and the first job
// failure.
func Run(ctx context.Context, jobs []Job) (map[string]Status, error) {
	if err := validate(jobs); err != nil {
		return nil, err
	}
	nodes := make(map[string]*node, len(jobs))
	for _, j := range jobs {
		nodes[j.Name] = &node{job: j, done: make(chan struct{})}
	}

	var mu sync.Mutex
	var firstErr error
	ctxSkipped := false

	var wg sync.WaitGroup
	for _, j := range jobs {
		n := nodes[j.Name]
		wg.Go(func() {
			defer close(n.done)
			for _, d := range n.job.Deps {
				dep := nodes[d]
				select {
				case <-dep.done:
				case <-ctx.Done():
				}
				if ctx.Err() != nil {
					mu.Lock()
					ctxSkipped = true
					mu.Unlock()
					return
				}
				if dep.status != Done {
					return // Skipped
				}
			}
			if ctx.Err() != nil {
				mu.Lock()
				ctxSkipped = true
				mu.Unlock()
				return
			}
			if err := n.job.Run(ctx); err != nil {
				n.status = Failed
				mu.Lock()
				if firstErr == nil {
					firstErr = fmt.Errorf("job %q: %w", n.job.Name, err)
				}
				mu.Unlock()
				return
			}
			n.status = Done
		})
	}
	wg.Wait()

	report := make(map[string]Status, len(nodes))
	for name, n := range nodes {
		report[name] = n.status
	}
	if firstErr != nil {
		return report, firstErr
	}
	if ctxSkipped {
		return report, ctx.Err()
	}
	return report, nil
}

// validate rejects duplicate names, unknown dependencies and cycles.
func validate(jobs []Job) error {
	byName := make(map[string]Job, len(jobs))
	for _, j := range jobs {
		if _, dup := byName[j.Name]; dup {
			return fmt.Errorf("%w: %q", ErrDuplicate, j.Name)
		}
		byName[j.Name] = j
	}
	for _, j := range jobs {
		for _, d := range j.Deps {
			if _, ok := byName[d]; !ok {
				return fmt.Errorf("%w: %q needs %q", ErrUnknownDep, j.Name, d)
			}
		}
	}

	// Kahn's algorithm: repeatedly remove jobs whose deps are all removed.
	// Whatever can never be removed sits on (or behind) a cycle.
	waiting := make(map[string]int, len(jobs))
	dependents := map[string][]string{}
	var ready []string
	for _, j := range jobs {
		waiting[j.Name] = len(j.Deps)
		for _, d := range j.Deps {
			dependents[d] = append(dependents[d], j.Name)
		}
		if len(j.Deps) == 0 {
			ready = append(ready, j.Name)
		}
	}
	for len(ready) > 0 {
		name := ready[0]
		ready = ready[1:]
		for _, dep := range dependents[name] {
			waiting[dep]--
			if waiting[dep] == 0 {
				ready = append(ready, dep)
			}
		}
	}
	for _, j := range jobs {
		if waiting[j.Name] > 0 {
			return fmt.Errorf("%w: involving %q", ErrCycle, j.Name)
		}
	}
	return nil
}
