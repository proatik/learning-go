# Deadline Budget

A request has one deadline but several sequential steps (auth, load, render…). Split the remaining time between the steps by **weight**, give each step its own timeout, and when one runs out, report **which** step it was via the timeout's cause.

```go
var ErrNoDeadline = errors.New("budget: ctx has no deadline")

type Step struct {
	Name   string
	Weight int
	Run    func(ctx context.Context) error
}

type StepTimeoutError struct {
	Step   string
	Budget time.Duration
}

func (e *StepTimeoutError) Error() string // `step "load" timed out after 4s`

func RunSteps(ctx context.Context, steps []Step) error
```

## Examples

Deadline 10s from the call, steps `A` (weight 1), `B` (2), `C` (2): the budgets are 2s, 4s and 4s.

| Steps take | Returns | After | Step deadlines |
|---|---|---|---|
| A 1s, B 3s, C 1s | `nil` | 5s | A at 2s, B at 1s+4s = 5s, C at 4s+4s = 8s |
| A 1s, B 5s | `&StepTimeoutError{Step: "B", Budget: 4s}` (C never runs) | 5s | |
| A fails with `boom` at 0.5s | `step "A": boom` (B, C never run) | 0.5s | |
| B is running when `ctx` is canceled at 1.5s | `context.Canceled` (not a `StepTimeoutError`) | 1.5s | |
| any steps, `ctx` without a deadline | `ErrNoDeadline` (no step runs) | 0s | |
| no steps, `ctx` with a deadline | `nil` | 0s | |

## Rules
- No deadline on `ctx` → return `ErrNoDeadline` without running anything.
- At the start, `total := time.Until(deadline)`. A weight `< 1` counts as 1. Step `i` gets `budget = total * weight / sum(weights)` (as `time.Duration` arithmetic, multiply first). Budgets are fixed up front: time a step doesn't use is **not** passed on.
- Run the steps in order. Before each one, if `ctx` is already done, return `context.Cause(ctx)`.
- Run each step with `context.WithTimeoutCause(ctx, budget, &StepTimeoutError{Step: name, Budget: budget})`, and cancel that context once the step returns.
- A step that returns `nil` succeeded, even if it ran past its budget. Move on to the next.
- A step that returns an error stops everything:
  - if **its** context is done, return `context.Cause` of it: your `*StepTimeoutError` when the step's budget ran out, or the parent's cause (e.g. `context.Canceled`) when the parent ended;
  - otherwise return `fmt.Errorf("step %q: %w", name, err)`.
- `StepTimeoutError.Error()` formats as `step "<name>" timed out after <budget>` (`%q` and `%v`). Callers find it with `errors.As`.
