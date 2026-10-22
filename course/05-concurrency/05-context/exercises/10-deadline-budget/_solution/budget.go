package budget

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// ErrNoDeadline is returned when the context has no deadline to split.
var ErrNoDeadline = errors.New("budget: ctx has no deadline")

// Step is one stage of a request. Weight decides its share of the deadline.
type Step struct {
	Name   string
	Weight int
	Run    func(ctx context.Context) error
}

// StepTimeoutError is the cancellation cause when a step uses up its budget.
type StepTimeoutError struct {
	Step   string
	Budget time.Duration
}

// Error formats e as `step "<name>" timed out after <budget>`.
func (e *StepTimeoutError) Error() string {
	return fmt.Sprintf("step %q timed out after %v", e.Step, e.Budget)
}

// RunSteps runs steps in order, giving each a timeout proportional to its
// weight out of the time left until ctx's deadline. It reports a step that
// ran out of time with a *StepTimeoutError cause.
func RunSteps(ctx context.Context, steps []Step) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrNoDeadline
	}
	total := time.Until(deadline)
	sum := 0
	for _, s := range steps {
		sum += weight(s)
	}
	for _, s := range steps {
		if ctx.Err() != nil {
			return context.Cause(ctx)
		}
		budget := total * time.Duration(weight(s)) / time.Duration(sum)
		if err := runStep(ctx, s, budget); err != nil {
			return err
		}
	}
	return nil
}

// runStep runs s with its own timeout and translates its failure.
func runStep(ctx context.Context, s Step, budget time.Duration) error {
	stepCtx, cancel := context.WithTimeoutCause(ctx, budget, &StepTimeoutError{Step: s.Name, Budget: budget})
	defer cancel()
	err := s.Run(stepCtx)
	if err == nil {
		return nil
	}
	// Decide before the deferred cancel runs: afterwards stepCtx is always done.
	if stepCtx.Err() != nil {
		return context.Cause(stepCtx)
	}
	return fmt.Errorf("step %q: %w", s.Name, err)
}

// weight returns s.Weight, counting anything below 1 as 1.
func weight(s Step) int {
	if s.Weight < 1 {
		return 1
	}
	return s.Weight
}
