// Package grades summarises exam scores.
package grades

import (
	"fmt"

	"learngo/course/04-generics-packages-testing/04-packages/exercises/01-split-package/stats"
)

// Report summarises scores as "n=<count> mean=<mean> min=<min> max=<max>",
// or "n=0" when there are no scores.
func Report(scores []int) string {
	lo, hi, ok := stats.MinMax(scores)
	if !ok {
		return "n=0"
	}
	return fmt.Sprintf("n=%d mean=%.2f min=%d max=%d", len(scores), stats.Mean(scores), lo, hi)
}
