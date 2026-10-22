package debounce

import "time"

// Debounce returns a channel that receives the latest value from in once in
// has been quiet for d. When in is closed, the pending value (if any) is
// sent and the output is closed.
func Debounce[T any](in <-chan T, d time.Duration) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		timer := time.NewTimer(d)
		timer.Stop()
		var (
			pending T
			fire    <-chan time.Time // nil while nothing is pending
		)
		for {
			select {
			case v, ok := <-in:
				if !ok {
					if fire != nil {
						out <- pending
					}
					return
				}
				pending = v
				timer.Reset(d)
				fire = timer.C
			case <-fire:
				out <- pending
				fire = nil
			}
		}
	}()
	return out
}
