package sales

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

// Parse converts "region,product,qty,price" lines into records.
func Parse(lines []string) ([]Record, error) {
	// TODO
	return nil, nil
}

// Summarize aggregates records per region, highest revenue first.
func Summarize(recs []Record) []RegionSummary {
	// TODO
	return nil
}

// ToJSON encodes the summaries as a JSON array.
func ToJSON(s []RegionSummary) (string, error) {
	// TODO
	return "", nil
}

// FromJSON decodes a JSON array produced by ToJSON.
func FromJSON(data string) ([]RegionSummary, error) {
	// TODO
	return nil, nil
}
