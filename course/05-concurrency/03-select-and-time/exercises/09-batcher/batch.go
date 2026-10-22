package batch

import "time"

// Batch groups values from in into slices of up to size items. A batch is
// sent when it is full or when wait has passed since its first item, and the
// partial batch is sent when in is closed.
func Batch[T any](in <-chan T, size int, wait time.Duration) <-chan []T {
	// TODO
	return nil
}
