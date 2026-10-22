package merge

// Merge returns a channel that receives every value from a and b and is
// closed once both are closed. A nil input counts as already closed.
func Merge[T any](a, b <-chan T) <-chan T {
	// TODO
	return nil
}
