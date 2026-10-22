// Package grades summarises exam scores.
package grades

import (
	"fmt"

	// The stats package lives in the stats/ sub-directory of this exercise.
	"learngo/course/04-generics-packages-testing/04-packages/exercises/01-split-package/stats"
)

// Report summarises scores as "n=<count> mean=<mean> min=<min> max=<max>",
// or "n=0" when there are no scores.
func Report(scores []int) string {
	_ = stats.Sum // TODO: call stats.Mean and stats.MinMax instead of the helpers below

	lo, hi, ok := minMax(scores)
	if !ok {
		return "n=0"
	}
	return fmt.Sprintf("n=%d mean=%.2f min=%d max=%d", len(scores), mean(scores), lo, hi)
}

// TODO: move the three helpers below into stats/stats.go (exported), then delete them here.

func sum(xs []int) int {
	total := 0
	for _, x := range xs {
		total += x
	}
	return total
}

func mean(xs []int) float64 {
	if len(xs) == 0 {
		return 0
	}
	return float64(sum(xs)) / float64(len(xs))
}

func minMax(xs []int) (lo, hi int, ok bool) {
	if len(xs) == 0 {
		return 0, 0, false
	}
	lo, hi = xs[0], xs[0]
	for _, x := range xs[1:] {
		lo = min(lo, x)
		hi = max(hi, x)
	}
	return lo, hi, true
}
