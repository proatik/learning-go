package heartbeat

import "time"

// Every calls f(1), f(2), … once per interval until stop is closed, then
// returns the number of calls. A non-positive interval returns 0 at once.
func Every(interval time.Duration, stop <-chan struct{}, f func(n int)) int {
	// TODO
	return 0
}
