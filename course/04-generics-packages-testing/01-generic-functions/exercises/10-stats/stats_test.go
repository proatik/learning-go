package stats_test

import (
	"math"
	"slices"
	"testing"

	st "learngo/course/04-generics-packages-testing/01-generic-functions/exercises/10-stats"
)

type (
	Celsius float64
	Cents   uint32
)

// near reports whether a and b are equal within a small tolerance.
func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

type result struct {
	name   string
	got    float64
	ok     bool
	want   float64
	wantOK bool
}

func check(t *testing.T, rs []result) {
	t.Helper()
	for _, r := range rs {
		if !near(r.got, r.want) || r.ok != r.wantOK {
			t.Errorf("%s = %v, %t, want %v, %t", r.name, r.got, r.ok, r.want, r.wantOK)
		}
	}
}

func r(name string, got float64, ok bool, want float64, wantOK bool) result {
	return result{name, got, ok, want, wantOK}
}

func TestMean(t *testing.T) {
	var rs []result
	g, ok := st.Mean([]int{1, 2, 3, 4})
	rs = append(rs, r("Mean([1 2 3 4])", g, ok, 2.5, true))
	g, ok = st.Mean([]int{-3, 3})
	rs = append(rs, r("Mean([-3 3])", g, ok, 0, true))
	g, ok = st.Mean([]int{7})
	rs = append(rs, r("Mean([7])", g, ok, 7, true))
	g, ok = st.Mean([]float64{0.5, 1.5, 4})
	rs = append(rs, r("Mean([0.5 1.5 4])", g, ok, 2, true))
	g, ok = st.Mean([]Celsius{20, 21, 22.5})
	rs = append(rs, r("Mean([]Celsius{20 21 22.5})", g, ok, 63.5/3, true))
	g, ok = st.Mean([]float32{1, 2})
	rs = append(rs, r("Mean([]float32{1 2})", g, ok, 1.5, true))
	g, ok = st.Mean([]int{})
	rs = append(rs, r("Mean([])", g, ok, 0, false))
	g, ok = st.Mean([]float64(nil))
	rs = append(rs, r("Mean(nil)", g, ok, 0, false))
	check(t, rs)
}

func TestMeanNoOverflow(t *testing.T) {
	var rs []result
	g, ok := st.Mean([]uint8{200, 200})
	rs = append(rs, r("Mean([]uint8{200 200})", g, ok, 200, true))
	g, ok = st.Mean([]int8{-128, -128, 127})
	rs = append(rs, r("Mean([]int8{-128 -128 127})", g, ok, -43, true))
	g, ok = st.Mean([]Cents{4_000_000_000, 4_000_000_000})
	rs = append(rs, r("Mean([]Cents{4e9 4e9})", g, ok, 4e9, true))
	g, ok = st.Mean([]int64{math.MaxInt64, math.MaxInt64})
	rs = append(rs, r("Mean([]int64{MaxInt64 MaxInt64})", g, ok, float64(math.MaxInt64), true))
	check(t, rs)
}

func TestVariance(t *testing.T) {
	var rs []result
	g, ok := st.Variance([]int{2, 4, 4, 4, 5, 5, 7, 9})
	rs = append(rs, r("Variance([2 4 4 4 5 5 7 9])", g, ok, 4, true))
	g, ok = st.Variance([]int{5, 5, 5})
	rs = append(rs, r("Variance([5 5 5])", g, ok, 0, true))
	g, ok = st.Variance([]int{42})
	rs = append(rs, r("Variance([42])", g, ok, 0, true))
	g, ok = st.Variance([]float64{1, 2, 3, 4})
	rs = append(rs, r("Variance([1 2 3 4])", g, ok, 1.25, true))
	g, ok = st.Variance([]uint{1, 3})
	rs = append(rs, r("Variance([]uint{1 3})", g, ok, 1, true))
	g, ok = st.Variance([]uint8{0, 255})
	rs = append(rs, r("Variance([]uint8{0 255})", g, ok, 16256.25, true))
	g, ok = st.Variance([]Celsius{-1.5, 1.5})
	rs = append(rs, r("Variance([]Celsius{-1.5 1.5})", g, ok, 2.25, true))
	g, ok = st.Variance([]int(nil))
	rs = append(rs, r("Variance(nil)", g, ok, 0, false))
	check(t, rs)
}

func TestMedian(t *testing.T) {
	var rs []result
	g, ok := st.Median([]int{3, 1, 2})
	rs = append(rs, r("Median([3 1 2])", g, ok, 2, true))
	g, ok = st.Median([]int{4, 1, 3, 2})
	rs = append(rs, r("Median([4 1 3 2])", g, ok, 2.5, true))
	g, ok = st.Median([]int{10, -10})
	rs = append(rs, r("Median([10 -10])", g, ok, 0, true))
	g, ok = st.Median([]int{5, 1, 5, 9, 5})
	rs = append(rs, r("Median([5 1 5 9 5])", g, ok, 5, true))
	g, ok = st.Median([]Celsius{21.5})
	rs = append(rs, r("Median([]Celsius{21.5})", g, ok, 21.5, true))
	g, ok = st.Median([]uint8{250, 255})
	rs = append(rs, r("Median([]uint8{250 255})", g, ok, 252.5, true))
	g, ok = st.Median([]float64{2.5, -1, 0.5, 100})
	rs = append(rs, r("Median([2.5 -1 0.5 100])", g, ok, 1.5, true))
	g, ok = st.Median([]int{})
	rs = append(rs, r("Median([])", g, ok, 0, false))
	check(t, rs)
}

func TestMedianDoesNotModify(t *testing.T) {
	in := []int{3, 1, 2, 5, 4}
	st.Median(in)
	st.Summarize(in)
	if want := []int{3, 1, 2, 5, 4}; !slices.Equal(in, want) {
		t.Errorf("Median(%v) modified the input: %v, want %v", want, in, want)
	}
}

func TestSummarize(t *testing.T) {
	got, ok := st.Summarize([]int{2, 4, 4, 4, 5, 5, 7, 9})
	want := st.Summary{Count: 8, Min: 2, Max: 9, Mean: 5, Median: 4.5, StdDev: 2}
	if !ok || got != want {
		t.Errorf("Summarize([2 4 4 4 5 5 7 9]) = %+v, %t, want %+v, %t", got, ok, want, true)
	}

	temps := []Celsius{19.5, 22, 20.5}
	got, ok = st.Summarize(temps)
	if !ok || got.Count != 3 || got.Min != 19.5 || got.Max != 22 || got.Median != 20.5 ||
		!near(got.Mean, 62.0/3) || !near(got.StdDev, math.Sqrt(1.0555555555555556)) {
		t.Errorf("Summarize(%v) = %+v, %t, want {Count:3 Min:19.5 Max:22 Mean:20.67 Median:20.5 StdDev:1.03}, true", temps, got, ok)
	}

	bytes := []uint8{255, 0, 255}
	got, ok = st.Summarize(bytes)
	if !ok || got.Count != 3 || got.Min != 0 || got.Max != 255 || got.Mean != 170 || got.Median != 255 {
		t.Errorf("Summarize(%v) = %+v, %t, want {Count:3 Min:0 Max:255 Mean:170 Median:255 ...}, true", bytes, got, ok)
	}

	got, ok = st.Summarize([]int{-7})
	want = st.Summary{Count: 1, Min: -7, Max: -7, Mean: -7, Median: -7, StdDev: 0}
	if !ok || got != want {
		t.Errorf("Summarize([-7]) = %+v, %t, want %+v, %t", got, ok, want, true)
	}

	for _, in := range [][]float64{nil, {}} {
		if got, ok := st.Summarize(in); ok || got != (st.Summary{}) {
			t.Errorf("Summarize(%v) = %+v, %t, want %+v, %t", in, got, ok, st.Summary{}, false)
		}
	}
}
