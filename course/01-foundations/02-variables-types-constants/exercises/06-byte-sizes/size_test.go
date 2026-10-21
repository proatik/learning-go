package size

import "testing"

func TestConstants(t *testing.T) {
	tests := []struct {
		name      string
		got, want int64
	}{
		{"KB", KB, 1024},
		{"MB", MB, 1024 * 1024},
		{"GB", GB, 1024 * 1024 * 1024},
		{"TB", TB, 1024 * 1024 * 1024 * 1024},
	}
	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s = %d, want %d", tt.name, tt.got, tt.want)
		}
	}
}

func TestHumanSize(t *testing.T) {
	const (
		kb = 1 << 10
		mb = 1 << 20
		gb = 1 << 30
		tb = 1 << 40
	)
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1023, "1023 B"},
		{kb, "1.0 KB"},
		{1536, "1.5 KB"},
		{mb - 1, "1024.0 KB"},
		{mb, "1.0 MB"},
		{5 * mb, "5.0 MB"},
		{gb - 1, "1024.0 MB"},
		{3*gb + 512*mb, "3.5 GB"},
		{tb, "1.0 TB"},
		{2 * tb, "2.0 TB"},
		{2048 * tb, "2048.0 TB"},
	}
	for _, tt := range tests {
		if got := HumanSize(tt.n); got != tt.want {
			t.Errorf("HumanSize(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}
