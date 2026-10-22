package timed

import "time"

// Timed runs every task concurrently and returns the wall time of the whole
// batch and the time each task took, in input order.
func Timed(tasks []func()) (total time.Duration, each []time.Duration) {
	// TODO
	return 0, nil
}
