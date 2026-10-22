package retry

import (
	"context"
	"errors"
	"time"
)

// ErrNoTime reports that the context's deadline would pass before the next
// attempt could start.
var ErrNoTime = errors.New("retry: not enough time for another attempt")

// Retry calls op up to attempts times, waiting base, 2*base, 4*base, …
// between failures. It stops early when ctx is done, or when the next wait
// would end at or after ctx's deadline.
func Retry(ctx context.Context, attempts int, base time.Duration, op func(context.Context) error) error {
	// TODO
	return nil
}
