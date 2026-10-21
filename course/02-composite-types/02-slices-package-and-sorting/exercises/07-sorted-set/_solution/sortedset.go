package sortedset

import "slices"

// Has reports whether v is in the sorted slice set.
func Has(set []int, v int) bool {
	_, found := slices.BinarySearch(set, v)
	return found
}

// InsertSorted inserts v into the sorted slice s at its sorted position,
// even if v is already present, and returns the updated slice.
func InsertSorted(s []int, v int) []int {
	// BinarySearch returns the insertion point whether or not v is found.
	i, _ := slices.BinarySearch(s, v)
	return slices.Insert(s, i, v)
}

// Add inserts v into the sorted slice set if it is not already present.
// It returns the updated slice and whether v was added.
func Add(set []int, v int) ([]int, bool) {
	i, found := slices.BinarySearch(set, v)
	if found {
		return set, false
	}
	return slices.Insert(set, i, v), true
}

// Remove deletes one occurrence of v from the sorted slice set.
// It returns the updated slice and whether v was found.
func Remove(set []int, v int) ([]int, bool) {
	i, found := slices.BinarySearch(set, v)
	if !found {
		return set, false
	}
	return slices.Delete(set, i, i+1), true // removes exactly set[i]
}
