// Package report summarises orders using the lazy adapters in package pipe.
package report

import (
	"fmt"
	"iter"
	"slices"

	"learngo/course/04-generics-packages-testing/06-module-challenge/exercises/06-lazy-pipeline/pipe"
)

// Order is one customer order. Cents is never negative.
type Order struct {
	ID       int
	Customer string
	Cents    int64
	Refunded bool
}

// paid reports whether o was not refunded.
func paid(o Order) bool { return !o.Refunded }

// Revenue returns the total Cents of all orders that were not refunded.
func Revenue(orders iter.Seq[Order]) int64 {
	return pipe.Fold(pipe.Filter(orders, paid), 0, func(sum int64, o Order) int64 {
		return sum + o.Cents
	})
}

// LargeOrderIDs returns the IDs of the first n paid orders of at least
// minCents, in input order. It stops reading orders once it has n.
func LargeOrderIDs(orders iter.Seq[Order], minCents int64, n int) []int {
	large := pipe.Filter(orders, func(o Order) bool { return paid(o) && o.Cents >= minCents })
	ids := slices.Collect(pipe.Take(pipe.Map(large, func(o Order) int { return o.ID }), n))
	if ids == nil {
		ids = []int{}
	}
	return ids
}

// Line formats o as "#ID Customer $D.CC", adding " (refunded)" if it was.
func Line(o Order) string {
	s := fmt.Sprintf("#%d %s $%d.%02d", o.ID, o.Customer, o.Cents/100, o.Cents%100)
	if o.Refunded {
		s += " (refunded)"
	}
	return s
}

// Pages yields the orders formatted with Line, size lines per page.
func Pages(orders iter.Seq[Order], size int) iter.Seq[[]string] {
	return pipe.Batch(pipe.Map(orders, Line), size)
}
