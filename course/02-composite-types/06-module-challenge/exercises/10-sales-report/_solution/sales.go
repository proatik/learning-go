package sales

import (
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

// Record is one parsed sales line.
type Record struct {
	Region  string
	Product string
	Qty     int
	Cents   int64 // unit price in cents
}

// RegionSummary aggregates the records of one region.
type RegionSummary struct {
	Region   string `json:"region"`
	Orders   int    `json:"orders"`
	Units    int    `json:"units"`
	Revenue  int64  `json:"revenue_cents"`
	AvgOrder int64  `json:"avg_order_cents"`
	Top      string `json:"top_product"`
}

// allDigits reports whether s is non-empty and only ASCII digits.
func allDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := range len(s) {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// parsePrice converts "D", "D.F" or "D.FF" into cents.
func parsePrice(s string) (int64, error) {
	whole, frac, hasDot := strings.Cut(s, ".")
	if !allDigits(whole) || (hasDot && (!allDigits(frac) || len(frac) > 2)) {
		return 0, fmt.Errorf("invalid price %q", s)
	}
	units, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price %q: %w", s, err)
	}
	var cents int64
	if hasDot {
		cents, _ = strconv.ParseInt(frac, 10, 64) // 1-2 digits: cannot fail
		if len(frac) == 1 {
			cents *= 10 // "2.5" means 50 cents
		}
	}
	return units*100 + cents, nil
}

// parseLine converts one data line into a Record.
func parseLine(line string) (Record, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		return Record{}, fmt.Errorf("got %d fields, want 4", len(fields))
	}
	for i := range fields {
		fields[i] = strings.TrimSpace(fields[i])
	}
	r := Record{Region: fields[0], Product: fields[1]}
	if r.Region == "" || r.Product == "" {
		return Record{}, errors.New("empty region or product")
	}
	qty, err := strconv.Atoi(fields[2])
	if err != nil || qty <= 0 {
		return Record{}, fmt.Errorf("invalid quantity %q", fields[2])
	}
	r.Qty = qty
	if r.Cents, err = parsePrice(fields[3]); err != nil {
		return Record{}, err
	}
	return r, nil
}

// Parse converts "region,product,qty,price" lines into records.
func Parse(lines []string) ([]Record, error) {
	recs := []Record{}
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		r, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("sales: line %d: %w", i+1, err)
		}
		recs = append(recs, r)
	}
	return recs, nil
}

// Summarize aggregates records per region, highest revenue first.
func Summarize(recs []Record) []RegionSummary {
	byRegion := map[string]*RegionSummary{}
	productRevenue := map[string]map[string]int64{} // region -> product -> revenue
	for _, r := range recs {
		s, ok := byRegion[r.Region]
		if !ok {
			s = &RegionSummary{Region: r.Region}
			byRegion[r.Region] = s
			productRevenue[r.Region] = map[string]int64{}
		}
		rev := int64(r.Qty) * r.Cents
		s.Orders++
		s.Units += r.Qty
		s.Revenue += rev
		productRevenue[r.Region][r.Product] += rev
	}

	out := make([]RegionSummary, 0, len(byRegion))
	for region, s := range byRegion {
		s.AvgOrder = s.Revenue / int64(s.Orders)
		// Visit products in sorted order so ties resolve alphabetically.
		products := productRevenue[region]
		best := int64(-1)
		for _, p := range slices.Sorted(maps.Keys(products)) {
			if products[p] > best {
				best, s.Top = products[p], p
			}
		}
		out = append(out, *s)
	}
	slices.SortFunc(out, func(a, b RegionSummary) int {
		return cmp.Or(cmp.Compare(b.Revenue, a.Revenue), cmp.Compare(a.Region, b.Region))
	})
	return out
}

// ToJSON encodes the summaries as a JSON array.
func ToJSON(s []RegionSummary) (string, error) {
	if s == nil {
		s = []RegionSummary{} // a nil slice would encode as null
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// FromJSON decodes a JSON array produced by ToJSON.
func FromJSON(data string) ([]RegionSummary, error) {
	var s []RegionSummary
	if err := json.Unmarshal([]byte(data), &s); err != nil {
		return nil, fmt.Errorf("sales: %w", err)
	}
	return s, nil
}
