// Package report summarises orders using the lazy adapters in package pipe.
package report

import (
	"iter"

	"learngo/course/04-generics-packages-testing/06-module-challenge/exercises/06-lazy-pipeline/pipe"
)

// Order is one customer order. Cents is never negative.
type Order struct {
	ID       int
	Customer string
	Cents    int64
	Refunded bool
}

// Revenue returns the total Cents of all orders that were not refunded.
func Revenue(orders iter.Seq[Order]) int64 {
	// TODO: pipe.Fold over pipe.Filter.
	_ = pipe.Fold[Order, int64]
	return 0
}

// LargeOrderIDs returns the IDs of the first n paid orders of at least
// minCents, in input order. It stops reading orders once it has n.
func LargeOrderIDs(orders iter.Seq[Order], minCents int64, n int) []int {
	// TODO: Filter, Map to the ID, Take n, then collect.
	return nil
}

// Line formats o as "#ID Customer $D.CC", adding " (refunded)" if it was.
func Line(o Order) string {
	// TODO
	return ""
}

// Pages yields the orders formatted with Line, size lines per page.
func Pages(orders iter.Seq[Order], size int) iter.Seq[[]string] {
	// TODO: pipe.Batch over pipe.Map.
	return func(yield func([]string) bool) {}
}
