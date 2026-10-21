package gradebook

// Gradebook maps a student's name to their scores by subject.
type Gradebook map[string]map[string]int

// AddGrade records (or overwrites) student's score in subject. g must be
// non-nil; the student's inner map is created on first use.
func AddGrade(g Gradebook, student, subject string, score int) {
	// TODO
}

// Average returns the mean of student's scores. ok is false if the student
// is unknown or has no grades.
func Average(g Gradebook, student string) (avg float64, ok bool) {
	// TODO
	return 0, false
}

// SubjectAverages returns, for each subject, the mean score of the students
// who have a score in it. The result is never nil.
func SubjectAverages(g Gradebook) map[string]float64 {
	// TODO
	return nil
}

// PrintReport prints one line per student (sorted by name) listing their
// scores (sorted by subject) and their average.
func PrintReport(g Gradebook) {
	// TODO
}
