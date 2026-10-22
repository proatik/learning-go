package budget

import (
	"context"
	"errors"
	"slices"
	"testing"
	"testing/synctest"
	"time"
)

var errBoom = errors.New("boom")

const ms = time.Millisecond

type key struct{}

// recorder builds steps that take a set time (or fail) and remembers what
// each one saw.
type recorder struct {
	start     time.Time
	ran       []string
	deadlines []time.Duration // each step's ctx deadline, relative to start
	values    []any
}

// step returns a Step that runs for d (unless its ctx ends first, then it
// returns ctx.Err()), and then returns err.
func (r *recorder) step(name string, weight int, d time.Duration, err error) Step {
	return Step{Name: name, Weight: weight, Run: func(ctx context.Context) error {
		r.ran = append(r.ran, name)
		dl, _ := ctx.Deadline()
		r.deadlines = append(r.deadlines, dl.Sub(r.start))
		r.values = append(r.values, ctx.Value(key{}))
		select {
		case <-time.After(d):
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}}
}

func TestStepTimeoutErrorMessage(t *testing.T) {
	tests := []struct {
		e    *StepTimeoutError
		want string
	}{
		{&StepTimeoutError{Step: "load", Budget: 4 * time.Second}, `step "load" timed out after 4s`},
		{&StepTimeoutError{Step: "auth", Budget: 1500 * ms}, `step "auth" timed out after 1.5s`},
		{&StepTimeoutError{Step: "", Budget: 0}, `step "" timed out after 0s`},
	}
	for _, tt := range tests {
		if got := tt.e.Error(); got != tt.want {
			t.Errorf("(%+v).Error() = %q, want %q", *tt.e, got, tt.want)
		}
	}
}

func TestAllSucceed(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.WithValue(t.Context(), key{}, "v"), 10*time.Second)
		defer cancel()
		r := &recorder{start: time.Now()}
		steps := []Step{r.step("A", 1, time.Second, nil), r.step("B", 2, 3*time.Second, nil), r.step("C", 2, time.Second, nil)}
		err := RunSteps(ctx, steps)
		elapsed := time.Since(r.start)
		if err != nil {
			t.Fatalf("RunSteps(10s; A=1 takes 1s, B=2 takes 3s, C=2 takes 1s) = %v, want <nil>", err)
		}
		if want := []string{"A", "B", "C"}; !slices.Equal(r.ran, want) {
			t.Errorf("RunSteps(10s; A, B, C) ran %v, want %v", r.ran, want)
		}
		if want := []time.Duration{2 * time.Second, 5 * time.Second, 8 * time.Second}; !slices.Equal(r.deadlines, want) {
			t.Errorf("RunSteps(10s; weights 1, 2, 2): step deadlines (from start) = %v, want %v", r.deadlines, want)
		}
		if elapsed != 5*time.Second {
			t.Errorf("RunSteps(10s; A, B, C) returned after %v, want 5s", elapsed)
		}
		for i, v := range r.values {
			if v != "v" {
				t.Errorf("RunSteps: step %s ctx.Value(key) = %v, want \"v\" (derive step contexts from ctx)", r.ran[i], v)
			}
		}
	})
}

func TestStepTimesOut(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		r := &recorder{start: time.Now()}
		steps := []Step{r.step("A", 1, time.Second, nil), r.step("B", 2, 5*time.Second, nil), r.step("C", 2, time.Second, nil)}
		err := RunSteps(ctx, steps)
		elapsed := time.Since(r.start)
		var ste *StepTimeoutError
		if !errors.As(err, &ste) {
			t.Fatalf("RunSteps(10s; A 1s, B 5s with a 4s budget) = %v, want a *StepTimeoutError", err)
		}
		if ste.Step != "B" || ste.Budget != 4*time.Second {
			t.Errorf("RunSteps(10s; A 1s, B 5s) = %+v, want {Step:B Budget:4s}", *ste)
		}
		if err.Error() != `step "B" timed out after 4s` {
			t.Errorf("RunSteps(10s; A 1s, B 5s) error = %q, want %q", err.Error(), `step "B" timed out after 4s`)
		}
		if want := []string{"A", "B"}; !slices.Equal(r.ran, want) {
			t.Errorf("RunSteps(10s; A 1s, B 5s) ran %v, want %v (stop after the timeout)", r.ran, want)
		}
		if elapsed != 5*time.Second {
			t.Errorf("RunSteps(10s; A 1s, B 5s) returned after %v, want 5s", elapsed)
		}
	})
}

func TestFirstStepTimesOut(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 6*time.Second)
		defer cancel()
		r := &recorder{start: time.Now()}
		steps := []Step{r.step("auth", 1, time.Hour, nil), r.step("load", 1, 0, nil), r.step("render", 1, 0, nil)}
		err := RunSteps(ctx, steps)
		var ste *StepTimeoutError
		if !errors.As(err, &ste) || ste.Step != "auth" || ste.Budget != 2*time.Second {
			t.Fatalf("RunSteps(6s; auth takes 1h, weights 1, 1, 1) = %v, want %q", err, `step "auth" timed out after 2s`)
		}
		if got := time.Since(r.start); got != 2*time.Second {
			t.Errorf("RunSteps(6s; auth takes 1h) returned after %v, want 2s", got)
		}
	})
}

func TestStepFails(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		r := &recorder{start: time.Now()}
		steps := []Step{r.step("A", 1, 500*ms, errBoom), r.step("B", 2, time.Second, nil)}
		err := RunSteps(ctx, steps)
		if !errors.Is(err, errBoom) || err.Error() != `step "A": boom` {
			t.Fatalf("RunSteps(A fails with boom) = %v, want %q", err, `step "A": boom`)
		}
		if want := []string{"A"}; !slices.Equal(r.ran, want) {
			t.Errorf("RunSteps(A fails) ran %v, want %v", r.ran, want)
		}
		if got := time.Since(r.start); got != 500*ms {
			t.Errorf("RunSteps(A fails at 0.5s) returned after %v, want 500ms", got)
		}
	})
}

func TestParentCanceled(t *testing.T) {
	ErrShutdown := errors.New("shutdown")
	for _, cause := range []error{nil, ErrShutdown} {
		synctest.Test(t, func(t *testing.T) {
			parent, stop := context.WithCancelCause(t.Context())
			defer stop(nil)
			ctx, cancel := context.WithTimeout(parent, 10*time.Second)
			defer cancel()
			time.AfterFunc(1500*ms, func() { stop(cause) })
			r := &recorder{start: time.Now()}
			steps := []Step{r.step("A", 1, time.Second, nil), r.step("B", 2, 3*time.Second, nil), r.step("C", 2, time.Second, nil)}
			err := RunSteps(ctx, steps)
			want := cause
			if want == nil {
				want = context.Canceled
			}
			var ste *StepTimeoutError
			if !errors.Is(err, want) || errors.As(err, &ste) {
				t.Errorf("RunSteps(parent canceled with %v at 1.5s, during B) = %v, want %v", cause, err, want)
			}
			if got := time.Since(r.start); got != 1500*ms {
				t.Errorf("RunSteps(parent canceled at 1.5s) returned after %v, want 1.5s", got)
			}
			if wantRan := []string{"A", "B"}; !slices.Equal(r.ran, wantRan) {
				t.Errorf("RunSteps(parent canceled at 1.5s) ran %v, want %v", r.ran, wantRan)
			}
		})
	}
}

func TestStepIgnoresItsBudget(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Second)
		defer cancel()
		var ran []string
		slow := Step{Name: "stubborn", Weight: 1, Run: func(context.Context) error {
			ran = append(ran, "stubborn")
			time.Sleep(3 * time.Second) // ignores its 2s budget, then succeeds
			return nil
		}}
		next := Step{Name: "next", Weight: 4, Run: func(context.Context) error {
			ran = append(ran, "next")
			return nil
		}}
		if err := RunSteps(ctx, []Step{slow, next}); err != nil {
			t.Errorf("RunSteps(stubborn step returns nil after its budget) = %v, want <nil>", err)
		}
		if want := []string{"stubborn", "next"}; !slices.Equal(ran, want) {
			t.Errorf("RunSteps(stubborn, next) ran %v, want %v", ran, want)
		}
	})
}

func TestWeightsBelowOne(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 9*time.Second)
		defer cancel()
		r := &recorder{start: time.Now()}
		steps := []Step{r.step("x", 0, 0, nil), r.step("y", 1, 0, nil), r.step("z", -3, 0, nil)}
		if err := RunSteps(ctx, steps); err != nil {
			t.Fatalf("RunSteps(9s; weights 0, 1, -3) = %v, want <nil>", err)
		}
		if want := []time.Duration{3 * time.Second, 3 * time.Second, 3 * time.Second}; !slices.Equal(r.deadlines, want) {
			t.Errorf("RunSteps(9s; weights 0, 1, -3 count as 1): step deadlines = %v, want %v", r.deadlines, want)
		}
	})
}

func TestNoDeadline(t *testing.T) {
	r := &recorder{start: time.Now()}
	err := RunSteps(t.Context(), []Step{r.step("A", 1, 0, nil)})
	if !errors.Is(err, ErrNoDeadline) {
		t.Errorf("RunSteps(ctx without deadline) = %v, want %v", err, ErrNoDeadline)
	}
	if len(r.ran) != 0 {
		t.Errorf("RunSteps(ctx without deadline) ran %v, want none", r.ran)
	}
}

func TestNoSteps(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		if err := RunSteps(ctx, nil); err != nil {
			t.Errorf("RunSteps(ctx with deadline, no steps) = %v, want <nil>", err)
		}
	})
}

func TestAlreadyExpired(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 0)
		defer cancel()
		r := &recorder{start: time.Now()}
		err := RunSteps(ctx, []Step{r.step("A", 1, time.Second, nil)})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("RunSteps(expired ctx) = %v, want %v", err, context.DeadlineExceeded)
		}
		if len(r.ran) != 0 {
			t.Errorf("RunSteps(expired ctx) ran %v, want none", r.ran)
		}
	})
}
