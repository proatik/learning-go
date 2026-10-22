package nonblock

// TrySend sends v on ch only if that can happen right now, and reports
// whether it did. It never blocks.
func TrySend[T any](ch chan<- T, v T) bool {
	select {
	case ch <- v:
		return true
	default:
		return false
	}
}

// TryRecv receives a value from ch only if one is ready right now. ok is
// false for an empty, closed or nil channel. It never blocks.
func TryRecv[T any](ch <-chan T) (v T, ok bool) {
	select {
	case v, ok = <-ch:
		return v, ok // ok is false if ch is closed
	default:
		var zero T
		return zero, false
	}
}

// DrainReady receives every value that is ready on ch right now, in order.
// It never blocks.
func DrainReady[T any](ch <-chan T) []T {
	var out []T
	for {
		v, ok := TryRecv(ch)
		if !ok {
			return out
		}
		out = append(out, v)
	}
}
