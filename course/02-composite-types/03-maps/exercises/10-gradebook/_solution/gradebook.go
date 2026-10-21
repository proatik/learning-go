package gradebook

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Gradebook maps a student's name to their scores by subject.
type Gradebook map[string]map[string]int

// AddGrade records (or overwrites) student's score in subject. g must be
// non-nil; the student's inner map is created on first use.
func AddGrade(g Gradebook, student, subject string, score int) {
	scores, ok := g[student]
	if !ok {
		scores = map[string]int{}
		g[student] = scores // store the handle: scores and g[student] are the same map
	}
	scores[subject] = score
}

// Average returns the mean of student's scores. ok is false if the student
// is unknown or has no grades.
func Average(g Gradebook, student string) (avg float64, ok bool) {
	scores := g[student] // nil for an unknown student; len(nil) == 0
	if len(scores) == 0 {
		return 0, false
	}
	sum := 0
	for _, s := range scores {
		sum += s
	}
	return float64(sum) / float64(len(scores)), true
}

// SubjectAverages returns, for each subject, the mean score of the students
// who have a score in it. The result is never nil.
func SubjectAverages(g Gradebook) map[string]float64 {
	sums := map[string]int{}
	counts := map[string]int{}
	for _, scores := range g {
		for subj, s := range scores {
			sums[subj] += s
			counts[subj]++
		}
	}
	avgs := make(map[string]float64, len(sums))
	for subj, sum := range sums {
		avgs[subj] = float64(sum) / float64(counts[subj])
	}
	return avgs
}

// PrintReport prints one line per student (sorted by name) listing their
// scores (sorted by subject) and their average.
func PrintReport(g Gradebook) {
	for _, name := range slices.Sorted(maps.Keys(g)) {
		avg, ok := Average(g, name)
		if !ok {
			fmt.Printf("%s: (no grades)\n", name)
			continue
		}
		scores := g[name]
		var parts []string
		for _, subj := range slices.Sorted(maps.Keys(scores)) {
			parts = append(parts, fmt.Sprintf("%s=%d", subj, scores[subj]))
		}
		fmt.Printf("%s: %s (avg %.2f)\n", name, strings.Join(parts, " "), avg)
	}
}
