package ranking

import (
	"cmp"
	"slices"
)

// Student is a named grade.
type Student struct {
	Name  string
	Grade int
}

// Ranked is a Student with its competition rank (1 = best).
type Ranked struct {
	Student
	Rank int
}

// compare orders by grade descending, then name ascending.
func compare(a, b Student) int {
	return cmp.Or(
		cmp.Compare(b.Grade, a.Grade), // b first: descending
		cmp.Compare(a.Name, b.Name),
	)
}

// SortStudents sorts students in place by grade descending, then name ascending.
func SortStudents(students []Student) {
	slices.SortFunc(students, compare)
}

// TopN returns the best n students in ranking order, without modifying students.
func TopN(students []Student, n int) []Student {
	if n <= 0 {
		return nil
	}
	sorted := slices.Clone(students) // sort a copy, never the caller's slice
	SortStudents(sorted)
	return sorted[:min(n, len(sorted))]
}

// Ranks returns every student in ranking order with its competition rank,
// without modifying students.
func Ranks(students []Student) []Ranked {
	sorted := slices.Clone(students)
	SortStudents(sorted)
	out := make([]Ranked, len(sorted))
	for i, s := range sorted {
		rank := i + 1
		if i > 0 && s.Grade == sorted[i-1].Grade {
			rank = out[i-1].Rank // tie: share the previous rank
		}
		out[i] = Ranked{Student: s, Rank: rank}
	}
	return out
}
