package twosum

// TwoSum returns indices i < j with nums[i]+nums[j] == target. Among all
// such pairs it returns the one with the smallest j, then the smallest i.
// ok is false (and i, j are 0) if there is no such pair.
func TwoSum(nums []int, target int) (i, j int, ok bool) {
	seen := make(map[int]int, len(nums)) // value → first index where it appeared
	for j, n := range nums {
		// Comma-ok is required: a stored index of 0 is a real answer.
		if i, found := seen[target-n]; found {
			return i, j, true
		}
		if _, dup := seen[n]; !dup {
			seen[n] = j // keep the earliest index so i is as small as possible
		}
	}
	return 0, 0, false
}
