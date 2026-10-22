package psum

// sumChunk adds up one chunk. Call it once per chunk; the tests replace it
// with a slow version to check that chunks are summed concurrently.
var sumChunk = func(part []int) int {
	total := 0
	for _, v := range part {
		total += v
	}
	return total
}

// Sum splits nums into k chunks, sums each chunk in its own goroutine and
// returns the grand total and the per-chunk sums in chunk order.
func Sum(nums []int, k int) (total int, partials []int) {
	// TODO
	return 0, nil
}
