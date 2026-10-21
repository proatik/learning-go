package grade

// Grade converts a score from 0 to 100 into a letter grade "A" to "F".
// Scores outside that range return "invalid".
func Grade(score int) string {
	// A tagless switch picks the first case whose condition is true,
	// so later cases may rely on earlier ones having failed.
	switch {
	case score < 0 || score > 100:
		return "invalid"
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}
