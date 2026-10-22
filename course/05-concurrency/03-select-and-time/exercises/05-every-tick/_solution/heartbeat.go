package heartbeat

import "time"

// Every calls f(1), f(2), … once per interval until stop is closed, then
// returns the number of calls. A non-positive interval returns 0 at once.
func Every(interval time.Duration, stop <-chan struct{}, f func(n int)) int {
	if interval <= 0 {
		return 0
	}
	t := time.NewTicker(interval)
	defer t.Stop()
	n := 0
	for {
		select {
		case <-stop:
			return n
		case <-t.C:
			n++
			f(n)
		}
	}
}
