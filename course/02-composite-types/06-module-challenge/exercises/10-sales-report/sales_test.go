package sales

import (
	"slices"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	lines := []string{
		"# region,product,qty,price",
		"north, tea, 3, 2.50",
		"south,coffee,1,4",
		"",
		"   ",
		"  # indented comment",
		"north,coffee,2,4.05",
		"  east , green tea ,10, 0.99  ",
		"west,cake,1,0.5",
		"west,water,7,0",
		"west,gold,1,1000000",
	}
	want := []Record{
		{"north", "tea", 3, 250},
		{"south", "coffee", 1, 400},
		{"north", "coffee", 2, 405},
		{"east", "green tea", 10, 99},
		{"west", "cake", 1, 50},
		{"west", "water", 7, 0},
		{"west", "gold", 1, 100000000},
	}
	got, err := Parse(lines)
	if err != nil {
		t.Fatalf("Parse() error = %v, want nil", err)
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Parse() =\n %v\nwant\n %v", got, want)
	}
}

func TestParseNoData(t *testing.T) {
	for _, lines := range [][]string{nil, {}, {"", "# only a comment", "  "}} {
		got, err := Parse(lines)
		if err != nil || got == nil || len(got) != 0 {
			t.Errorf("Parse(%q) = %#v, %v, want empty slice, nil", lines, got, err)
		}
	}
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		line string
	}{
		{"north,tea,3"},
		{"north,tea,3,2.50,extra"},
		{"north"},
		{",tea,3,2.50"},
		{"north, ,3,2.50"},
		{"north,tea,0,2.50"},
		{"north,tea,-2,2.50"},
		{"north,tea,three,2.50"},
		{"north,tea,1.5,2.50"},
		{"north,tea,3,"},
		{"north,tea,3,-1"},
		{"north,tea,3,1."},
		{"north,tea,3,.5"},
		{"north,tea,3,1.234"},
		{"north,tea,3,abc"},
		{"north,tea,3,1.x"},
		{"north,tea,3,+4"},
		{"north,tea,3,1 2"},
		{"north,tea,3,1.2.3"},
	}
	for _, tt := range tests {
		lines := []string{"# header", "south,coffee,1,4", "", tt.line, "east,tea,1,1"}
		got, err := Parse(lines)
		if err == nil {
			t.Errorf("Parse(.., %q, ..) error = nil, want error", tt.line)
			continue
		}
		if got != nil {
			t.Errorf("Parse(.., %q, ..) = %v, want nil on error", tt.line, got)
		}
		if !strings.Contains(err.Error(), "line 4") {
			t.Errorf("Parse(.., %q, ..) error = %q, want it to mention \"line 4\"", tt.line, err)
		}
	}
	// "north,tea,3,1,5" has five fields: a comma is not a decimal point.
	if _, err := Parse([]string{"north,tea,3,1,5"}); err == nil || !strings.Contains(err.Error(), "line 1") {
		t.Errorf("Parse(\"north,tea,3,1,5\") error = %v, want an error mentioning line 1", err)
	}
}

func TestSummarize(t *testing.T) {
	recs := []Record{
		{"north", "tea", 3, 250},
		{"south", "coffee", 1, 400},
		{"north", "coffee", 2, 405},
	}
	want := []RegionSummary{
		{"north", 2, 5, 1560, 780, "coffee"},
		{"south", 1, 1, 400, 400, "coffee"},
	}
	if got := Summarize(recs); !slices.Equal(got, want) {
		t.Fatalf("Summarize(example) =\n %v\nwant\n %v", got, want)
	}
}

func TestSummarizeTiesAndOrder(t *testing.T) {
	recs := []Record{
		{"b", "zeta", 1, 100},
		{"b", "alpha", 2, 50}, // ties zeta on revenue → alpha wins
		{"a", "x", 1, 200},    // a and b tie on revenue → a first
		{"c", "big", 1, 999},
		{"d", "free", 5, 0},
		{"e", "pen", 1, 100},
		{"e", "pen", 1, 100}, // same product twice
		{"e", "ink", 3, 10},
		{"e", "pen", 1, 1},
	}
	want := []RegionSummary{
		{"c", 1, 1, 999, 999, "big"},
		{"e", 4, 6, 231, 57, "pen"}, // 231/4 truncates to 57
		{"a", 1, 1, 200, 200, "x"},
		{"b", 2, 3, 200, 100, "alpha"},
		{"d", 1, 5, 0, 0, "free"},
	}
	got := Summarize(recs)
	if !slices.Equal(got, want) {
		t.Fatalf("Summarize =\n %v\nwant\n %v", got, want)
	}
	// Running it again must give the same (deterministic) result.
	for range 20 {
		if again := Summarize(recs); !slices.Equal(again, want) {
			t.Fatalf("Summarize is not deterministic: %v", again)
		}
	}
}

func TestSummarizeEmpty(t *testing.T) {
	if got := Summarize(nil); len(got) != 0 {
		t.Errorf("Summarize(nil) = %v, want empty", got)
	}
}

func TestParseThenSummarize(t *testing.T) {
	recs, err := Parse([]string{
		"# region,product,qty,price",
		"north, tea, 3, 2.50",
		"south,coffee,1,4",
		"",
		"north,coffee,2,4.05",
	})
	if err != nil {
		t.Fatalf("Parse error = %v", err)
	}
	got := Summarize(recs)
	want := []RegionSummary{
		{"north", 2, 5, 1560, 780, "coffee"},
		{"south", 1, 1, 400, 400, "coffee"},
	}
	if !slices.Equal(got, want) {
		t.Fatalf("Summarize(Parse(example)) = %v, want %v", got, want)
	}
}

func TestToJSON(t *testing.T) {
	tests := []struct {
		in   []RegionSummary
		want string
	}{
		{
			[]RegionSummary{{"south", 1, 1, 400, 400, "coffee"}},
			`[{"region":"south","orders":1,"units":1,"revenue_cents":400,"avg_order_cents":400,"top_product":"coffee"}]`,
		},
		{
			[]RegionSummary{{"a", 2, 3, 10, 5, "x"}, {"b", 0, 0, 0, 0, ""}},
			`[{"region":"a","orders":2,"units":3,"revenue_cents":10,"avg_order_cents":5,"top_product":"x"},` +
				`{"region":"b","orders":0,"units":0,"revenue_cents":0,"avg_order_cents":0,"top_product":""}]`,
		},
		{nil, `[]`},
		{[]RegionSummary{}, `[]`},
	}
	for _, tt := range tests {
		got, err := ToJSON(tt.in)
		if err != nil || got != tt.want {
			t.Errorf("ToJSON(%v) = %s, %v, want %s, nil", tt.in, got, err, tt.want)
		}
	}
}

func TestJSONRoundTrip(t *testing.T) {
	in := []RegionSummary{
		{"north", 2, 5, 1560, 780, "coffee"},
		{"south \"quoted\"", 1, 1, 400, 400, "café"},
	}
	s, err := ToJSON(in)
	if err != nil {
		t.Fatalf("ToJSON error = %v", err)
	}
	out, err := FromJSON(s)
	if err != nil {
		t.Fatalf("FromJSON(%s) error = %v", s, err)
	}
	if !slices.Equal(out, in) {
		t.Fatalf("round trip = %v, want %v", out, in)
	}

	decoded, err := FromJSON(`[{"region":"x","orders":3,"units":4,"revenue_cents":50,"avg_order_cents":16,"top_product":"y"}]`)
	if want := []RegionSummary{{"x", 3, 4, 50, 16, "y"}}; err != nil || !slices.Equal(decoded, want) {
		t.Fatalf("FromJSON(tagged input) = %v, %v, want %v, nil (use the json tags)", decoded, err, want)
	}
}

func TestFromJSONInvalid(t *testing.T) {
	for _, data := range []string{"", "{", "not json", `[{"orders":"many"}]`, `{"region":"x"}`} {
		if got, err := FromJSON(data); err == nil || got != nil {
			t.Errorf("FromJSON(%q) = %v, %v, want nil, error", data, got, err)
		}
	}
}
