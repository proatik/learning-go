package duration

import "testing"

func TestHumanDuration(t *testing.T) {
	tests := []struct {
		seconds int
		want    string
	}{
		{0, "0s"},
		{1, "1s"},
		{59, "59s"},
		{60, "1m"},
		{61, "1m 1s"},
		{3599, "59m 59s"},
		{3600, "1h"},
		{3605, "1h 5s"},
		{3660, "1h 1m"},
		{3661, "1h 1m 1s"},
		{86399, "23h 59m 59s"},
		{86400, "1d"},
		{86401, "1d 1s"},
		{90000, "1d 1h"},
		{90061, "1d 1h 1m 1s"},
		{1_000_000, "11d 13h 46m 40s"},
		{100 * 86400, "100d"},
		{-1, "-1s"},
		{-90, "-1m 30s"},
		{-3600, "-1h"},
		{-90061, "-1d 1h 1m 1s"},
	}
	for _, tt := range tests {
		if got := HumanDuration(tt.seconds); got != tt.want {
			t.Errorf("HumanDuration(%d) = %q, want %q", tt.seconds, got, tt.want)
		}
	}
}
