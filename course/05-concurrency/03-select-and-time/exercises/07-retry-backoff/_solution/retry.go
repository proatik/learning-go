package retry

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrNoAttempts is returned when attempts < 1.
	ErrNoAttempts = errors.New("retry: no attempts")
	// ErrStopped is wrapped in the error returned when stop is closed while
	// waiting between attempts.
	ErrStopped = errors.New("retry: stopped")
)

// Do calls f up to attempts times, waiting base, 2·base, 4·base, … (at most
// maxDelay) between failed calls. It returns nil on the first success, or an
// error wrapping the last failure.
func Do(stop <-chan struct{}, attempts int, base, maxDelay time.Duration, f func() error) error {
	if attempts < 1 {
		return ErrNoAttempts
	}
	delay := base
	var err error
	for call := 1; ; call++ {
		if err = f(); err == nil {
			return nil
		}
		if call == attempts {
			return fmt.Errorf("retry: %d attempts failed: %w", attempts, err)
		}
		wait := time.NewTimer(min(delay, maxDelay))
		select {
		case <-wait.C:
		case <-stop:
			wait.Stop()
			return fmt.Errorf("%w after %d attempts: %w", ErrStopped, call, err)
		}
		delay *= 2
	}
}
