package batch

import "time"

// Batch groups values from in into slices of up to size items. A batch is
// sent when it is full or when wait has passed since its first item, and the
// partial batch is sent when in is closed.
func Batch[T any](in <-chan T, size int, wait time.Duration) <-chan []T {
	size = max(size, 1)
	out := make(chan []T)
	go func() {
		defer close(out)
		timer := time.NewTimer(wait)
		timer.Stop()
		var (
			batch   []T
			expired <-chan time.Time // nil while the batch is empty
		)
		flush := func() {
			timer.Stop()
			expired = nil
			out <- batch
			batch = nil // a fresh slice next time: the receiver owns this one
		}
		for {
			select {
			case v, ok := <-in:
				if !ok {
					if len(batch) > 0 {
						flush()
					}
					return
				}
				if len(batch) == 0 {
					timer.Reset(wait)
					expired = timer.C
				}
				batch = append(batch, v)
				if len(batch) == size {
					flush()
				}
			case <-expired:
				flush()
			}
		}
	}()
	return out
}
