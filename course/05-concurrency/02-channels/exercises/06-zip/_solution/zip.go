package zip

// Pair holds one value from each zipped channel.
type Pair[X, Y any] struct {
	First  X
	Second Y
}

// Zip returns a channel of pairs: each round it receives one value from a,
// then one from b. It closes the result as soon as either input is closed.
func Zip[X, Y any](a <-chan X, b <-chan Y) <-chan Pair[X, Y] {
	out := make(chan Pair[X, Y])
	go func() {
		defer close(out)
		for {
			x, ok := <-a
			if !ok {
				return
			}
			y, ok := <-b
			if !ok {
				return
			}
			out <- Pair[X, Y]{First: x, Second: y}
		}
	}()
	return out
}
