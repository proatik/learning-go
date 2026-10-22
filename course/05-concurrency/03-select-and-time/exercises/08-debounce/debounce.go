package debounce

import "time"

// Debounce returns a channel that receives the latest value from in once in
// has been quiet for d. When in is closed, the pending value (if any) is
// sent and the output is closed.
func Debounce[T any](in <-chan T, d time.Duration) <-chan T {
	// TODO
	return nil
}
