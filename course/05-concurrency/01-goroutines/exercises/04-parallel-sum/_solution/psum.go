package psum

import "sync"

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
	n := len(nums)
	k = max(min(k, n), 1)
	if n == 0 {
		return 0, []int{}
	}
	partials = make([]int, k)
	var wg sync.WaitGroup
	for i := range k {
		chunk := nums[i*n/k : (i+1)*n/k]
		wg.Go(func() { partials[i] = sumChunk(chunk) })
	}
	wg.Wait()
	for _, p := range partials {
		total += p
	}
	return total, partials
}
