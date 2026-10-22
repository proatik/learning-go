package retry

import (
	"errors"
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
	// TODO
	return nil
}
