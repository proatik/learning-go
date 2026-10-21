# Sales Report

Parse CSV-like sales lines into structs, aggregate them per region, and round-trip the report through JSON using struct tags.

```go
type Record struct {
	Region  string
	Product string
	Qty     int
	Cents   int64 // unit price in cents
}

type RegionSummary struct {
	Region   string `json:"region"`
	Orders   int    `json:"orders"`
	Units    int    `json:"units"`
	Revenue  int64  `json:"revenue_cents"`
	AvgOrder int64  `json:"avg_order_cents"`
	Top      string `json:"top_product"`
}

func Parse(lines []string) ([]Record, error)
func Summarize(recs []Record) []RegionSummary
func ToJSON(s []RegionSummary) (string, error)
func FromJSON(data string) ([]RegionSummary, error)
```

## Example

```go
recs, _ := Parse([]string{
	"# region,product,qty,price",
	"north, tea, 3, 2.50",
	"south,coffee,1,4",
	"",
	"north,coffee,2,4.05",
})
// [{north tea 3 250} {south coffee 1 400} {north coffee 2 405}]

Summarize(recs)
// [{north 2 5 1560 780 coffee} {south 1 1 400 400 coffee}]

ToJSON([]RegionSummary{{"south", 1, 1, 400, 400, "coffee"}})
// `[{"region":"south","orders":1,"units":1,"revenue_cents":400,"avg_order_cents":400,"top_product":"coffee"}]`
```

## Rules

**Parse** — each line is `region,product,qty,price`.
- Skip lines that are empty after `strings.TrimSpace`, and lines whose trimmed text starts with `#`.
- Split on `,`; there must be exactly **4** fields. Trim spaces around every field.
- `region` and `product` must be non-empty. `qty` must parse with `strconv.Atoi` and be `> 0`.
- `price` is `D`, `D.F` or `D.FF` where `D` and `F` are ASCII digits (`D` at least one digit): `"4"` → 400, `"2.5"` → 250, `"4.05"` → 405, `"0.99"` → 99. Anything else (`"-1"`, `"1."`, `".5"`, `"1.234"`, `"1,5"`, `"abc"`, `""`) is invalid.
- On the first bad line return `nil` and an error whose message contains `line N` (`N` = 1-based index into `lines`).
- No data lines → empty slice, `nil` error.

**Summarize** — one summary per region.
- `Orders` = number of records, `Units` = sum of `Qty`, `Revenue` = sum of `Qty*Cents`, `AvgOrder` = `Revenue / Orders` (integer division).
- `Top` = the product with the highest revenue in that region; ties go to the alphabetically smallest product.
- Sort by `Revenue` descending, then `Region` ascending. No records → empty slice.

**ToJSON / FromJSON**
- `ToJSON` uses `encoding/json` (`json.Marshal`) with the tags above; a nil or empty slice encodes as `"[]"`.
- `FromJSON` decodes such a string back (`json.Unmarshal`); invalid JSON returns `nil` and an error.
