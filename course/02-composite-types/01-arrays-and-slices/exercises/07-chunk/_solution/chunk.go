package chunk

// Chunk splits s into consecutive sub-slices of at most size elements.
// Every chunk shares s's backing array and has cap == len.
// It returns nil if size <= 0.
func Chunk(s []int, size int) [][]int {
	if size <= 0 {
		return nil
	}
	n := (len(s) + size - 1) / size // number of chunks, rounded up
	chunks := make([][]int, 0, n)
	for lo := 0; lo < len(s); lo += size {
		hi := min(lo+size, len(s)) // the last chunk may be shorter
		// s[lo:hi:hi] caps the chunk at its own end, so an append on it
		// must reallocate instead of overwriting s[hi].
		chunks = append(chunks, s[lo:hi:hi])
	}
	return chunks
}
