package retry

import (
	"errors"
	"fmt"
	"testing"
)

// timeout is a value-receiver error type with its own Temporary method.
type timeout struct{ temp bool }

func (t timeout) Error() string   { return fmt.Sprintf("timeout(temp=%t)", t.temp) }
func (t timeout) Temporary() bool { return t.temp }

// permanentWrapper says "not temporary" even though it wraps a temporary error.
type permanentWrapper struct{ err error }

func (p *permanentWrapper) Error() string   { return "permanent: " + p.err.Error() }
func (p *permanentWrapper) Temporary() bool { return false }
func (p *permanentWrapper) Unwrap() error   { return p.err }

var (
	errBusy   = errors.New("busy")
	errBadKey = errors.New("bad key")
)

func temp(err error) error { return &TemporaryError{Err: err} }

func TestTemporaryError(t *testing.T) {
	e := &TemporaryError{Err: errBusy}
	if got := e.Error(); got != "busy" {
		t.Errorf("(&TemporaryError{busy}).Error() = %q, want %q", got, "busy")
	}
	if !e.Temporary() {
		t.Errorf("(&TemporaryError{busy}).Temporary() = false, want true")
	}
	if !errors.Is(e, errBusy) {
		t.Errorf("errors.Is(&TemporaryError{busy}, errBusy) = false, want true (Unwrap)")
	}
}

func TestIsTemporary(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"plain", errBusy, false},
		{"TemporaryError", temp(errBusy), true},
		{"wrapped TemporaryError", fmt.Errorf("a: %w", fmt.Errorf("b: %w", temp(errBusy))), true},
		{"%v cuts the chain", fmt.Errorf("a: %v", temp(errBusy)), false},
		{"value type true", timeout{temp: true}, true},
		{"value type false", timeout{temp: false}, false},
		{"wrapped value type", fmt.Errorf("dial: %w", timeout{temp: true}), true},
		{"first one decides", &permanentWrapper{err: temp(errBusy)}, false},
		{"inside Join", errors.Join(errBadKey, temp(errBusy)), true},
	}
	for _, tt := range tests {
		if got := IsTemporary(tt.err); got != tt.want {
			t.Errorf("IsTemporary(%v) = %t, want %t (%s)", tt.err, got, tt.want, tt.name)
		}
	}
}

// script returns an f that returns results[0], results[1], ... and counts calls.
func script(results ...error) (func() error, *int) {
	calls := 0
	return func() error {
		calls++
		if calls > len(results) {
			return results[len(results)-1]
		}
		return results[calls-1]
	}, &calls
}

func TestRetry(t *testing.T) {
	tests := []struct {
		name      string
		attempts  int
		results   []error
		wantErr   string // "" means nil
		wantCalls int
	}{
		{"first try", 3, []error{nil}, "", 1},
		{"third try", 3, []error{temp(errBusy), temp(errBusy), nil}, "", 3},
		{"gave up", 3, []error{temp(errBusy)}, "gave up after 3 attempts: busy", 3},
		{"single attempt temp", 1, []error{temp(errBusy)}, "gave up after 1 attempts: busy", 1},
		{"permanent at once", 5, []error{errBadKey}, "attempt 1: bad key", 1},
		{"permanent later", 5, []error{temp(errBusy), errBadKey}, "attempt 2: bad key", 2},
		{"value type temp", 2, []error{timeout{temp: true}}, "gave up after 2 attempts: timeout(temp=true)", 2},
		{"value type permanent", 4, []error{timeout{temp: false}}, "attempt 1: timeout(temp=false)", 1},
		{"last error is reported", 3, []error{temp(errors.New("e1")), temp(errors.New("e2")), temp(errors.New("e3"))}, "gave up after 3 attempts: e3", 3},
	}
	for _, tt := range tests {
		f, calls := script(tt.results...)
		err := Retry(tt.attempts, f)
		if tt.wantErr == "" {
			if err != nil {
				t.Errorf("Retry(%d, f) error = %v, want <nil> (%s)", tt.attempts, err, tt.name)
			}
		} else if err == nil || err.Error() != tt.wantErr {
			t.Errorf("Retry(%d, f) error = %v, want %v (%s)", tt.attempts, err, tt.wantErr, tt.name)
		}
		if *calls != tt.wantCalls {
			t.Errorf("Retry(%d, f) called f %d times, want %d (%s)", tt.attempts, *calls, tt.wantCalls, tt.name)
		}
	}
}

func TestRetryKeepsChain(t *testing.T) {
	f, _ := script(temp(errBusy))
	err := Retry(2, f)
	if !errors.Is(err, errBusy) {
		t.Errorf("Retry(2, f) error = %v, want %v", err, errBusy)
	}
	if _, ok := errors.AsType[*TemporaryError](err); !ok {
		t.Errorf("Retry(2, f) error = %v, want a wrapped *TemporaryError", err)
	}

	f, _ = script(fmt.Errorf("open: %w", errBadKey))
	err = Retry(2, f)
	if !errors.Is(err, errBadKey) {
		t.Errorf("Retry(2, f) error = %v, want %v", err, errBadKey)
	}
}

func TestRetryNoAttempts(t *testing.T) {
	for _, n := range []int{0, -1} {
		f, calls := script(nil)
		if err := Retry(n, f); err != ErrNoAttempts {
			t.Errorf("Retry(%d, f) error = %v, want %v", n, err, ErrNoAttempts)
		}
		if *calls != 0 {
			t.Errorf("Retry(%d, f) called f %d times, want 0", n, *calls)
		}
	}
}
