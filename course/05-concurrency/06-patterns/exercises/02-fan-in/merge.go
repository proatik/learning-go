package fanin

// Merge returns one channel that yields every value received from every
// input. It is closed once all inputs are closed and drained. Nil inputs are
// ignored.
func Merge[T any](inputs ...<-chan T) <-chan T {
	// TODO
	return nil
}
