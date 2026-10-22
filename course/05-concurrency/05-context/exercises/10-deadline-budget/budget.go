package budget

import (
	"context"
	"errors"
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
	// TODO
	return ""
}

// RunSteps runs steps in order, giving each a timeout proportional to its
// weight out of the time left until ctx's deadline. It reports a step that
// ran out of time with a *StepTimeoutError cause.
func RunSteps(ctx context.Context, steps []Step) error {
	// TODO
	return nil
}
