package report

import "fmt"

// Summary describes how many items of a kind there are, e.g. "apples: 3 items".
func Summary(kind string, count int) string {
	return fmt.Sprintf("%s: %d items", kind, count)
}

// Percent formats part/total as a percentage with one decimal, e.g. "25.0%".
func Percent(part, total int) string {
	p := float64(part) / float64(total) * 100
	return fmt.Sprintf("%.1f%%", p)
}
