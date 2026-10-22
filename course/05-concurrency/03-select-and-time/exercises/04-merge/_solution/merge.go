package merge

// Merge returns a channel that receives every value from a and b and is
// closed once both are closed. A nil input counts as already closed.
func Merge[T any](a, b <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		for a != nil || b != nil {
			select {
			case v, ok := <-a:
				if !ok {
					a = nil // a receive from nil blocks forever: this case is now off
					continue
				}
				out <- v
			case v, ok := <-b:
				if !ok {
					b = nil
					continue
				}
				out <- v
			}
		}
	}()
	return out
}
