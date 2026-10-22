package batch

// Batch returns a channel that receives the values from in grouped into
// slices of size, in order; the last slice may be shorter. It is closed after
// in is closed. Batch panics if size < 1.
func Batch(in <-chan string, size int) <-chan []string {
	if size < 1 {
		panic("batch: size must be at least 1")
	}
	out := make(chan []string)
	go func() {
		defer close(out)
		buf := make([]string, 0, size)
		for v := range in {
			buf = append(buf, v)
			if len(buf) == size {
				out <- buf
				buf = make([]string, 0, size)
			}
		}
		if len(buf) > 0 {
			out <- buf
		}
	}()
	return out
}
