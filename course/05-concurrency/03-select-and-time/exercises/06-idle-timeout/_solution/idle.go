package idle

import "time"

// UntilIdle receives events until none has arrived for d, or events is
// closed, and returns what it received.
func UntilIdle(events <-chan string, d time.Duration) []string {
	var got []string
	idle := time.NewTimer(d)
	defer idle.Stop()
	for {
		select {
		case e, ok := <-events:
			if !ok {
				return got
			}
			got = append(got, e)
			idle.Reset(d)
		case <-idle.C:
			return got
		}
	}
}
