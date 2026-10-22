package tee

import "context"

// Tee copies every value from in to both returned channels, in order. Both
// outputs are closed when in is closed or ctx is done.
func Tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {
	out1, out2 := make(chan T), make(chan T)
	go func() {
		defer close(out1)
		defer close(out2)
		for {
			var v T
			select {
			case <-ctx.Done():
				return
			case x, ok := <-in:
				if !ok {
					return
				}
				v = x
			}
			// Deliver v to both, in whichever order the readers are ready;
			// a nil channel disables its case once it has its copy.
			o1, o2 := out1, out2
			for o1 != nil || o2 != nil {
				select {
				case <-ctx.Done():
					return
				case o1 <- v:
					o1 = nil
				case o2 <- v:
					o2 = nil
				}
			}
		}
	}()
	return out1, out2
}
