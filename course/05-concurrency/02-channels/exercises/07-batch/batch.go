package batch

// Batch returns a channel that receives the values from in grouped into
// slices of size, in order; the last slice may be shorter. It is closed after
// in is closed. Batch panics if size < 1.
func Batch(in <-chan string, size int) <-chan []string {
	// TODO
	return nil
}
