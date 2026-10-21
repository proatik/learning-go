package ranking

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

// SortStudents sorts students in place by grade descending, then name ascending.
func SortStudents(students []Student) {
	// TODO
}

// TopN returns the best n students in ranking order, without modifying students.
func TopN(students []Student, n int) []Student {
	// TODO
	return nil
}

// Ranks returns every student in ranking order with its competition rank,
// without modifying students.
func Ranks(students []Student) []Ranked {
	// TODO
	return nil
}
