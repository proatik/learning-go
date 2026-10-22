package idle

import "time"

// UntilIdle receives events until none has arrived for d, or events is
// closed, and returns what it received.
func UntilIdle(events <-chan string, d time.Duration) []string {
	// TODO
	return nil
}
