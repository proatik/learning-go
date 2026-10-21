package counter

import "testing"

func TestZeroValue(t *testing.T) {
	var c Counter
	if got := c.Value(); got != 0 {
		t.Errorf("zero Counter Value() = %d, want 0", got)
	}
}

func TestSequences(t *testing.T) {
	type op struct {
		name string
		arg  int
	}
	tests := []struct {
		desc string
		ops  []op
		want int
	}{
		{"Inc x2", []op{{"inc", 0}, {"inc", 0}}, 2},
		{"Inc x3, Dec", []op{{"inc", 0}, {"inc", 0}, {"inc", 0}, {"dec", 0}}, 2},
		{"Dec at zero", []op{{"dec", 0}}, 0},
		{"Dec twice from 1", []op{{"inc", 0}, {"dec", 0}, {"dec", 0}}, 0},
		{"Add 5", []op{{"add", 5}}, 5},
		{"Add 5, Add -2", []op{{"add", 5}, {"add", -2}}, 3},
		{"Add 3, Add -10", []op{{"add", 3}, {"add", -10}}, 0},
		{"Add -1 at zero", []op{{"add", -1}}, 0},
		{"Add 0", []op{{"add", 0}}, 0},
		{"Add 7, Reset", []op{{"add", 7}, {"reset", 0}}, 0},
		{"Reset, Inc", []op{{"add", 7}, {"reset", 0}, {"inc", 0}}, 1},
	}
	for _, tt := range tests {
		var c Counter
		for _, o := range tt.ops {
			switch o.name {
			case "inc":
				c.Inc()
			case "dec":
				c.Dec()
			case "add":
				c.Add(o.arg)
			case "reset":
				c.Reset()
			}
		}
		if got := c.Value(); got != tt.want {
			t.Errorf("%s: Value() = %d, want %d", tt.desc, got, tt.want)
		}
	}
}

func TestThroughPointer(t *testing.T) {
	c := &Counter{}
	p := c
	p.Add(4)
	if got := c.Value(); got != 4 {
		t.Errorf("Add(4) through a copy of the pointer: Value() = %d, want 4", got)
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		words []string
		want  map[string]int
	}{
		{nil, map[string]int{}},
		{[]string{"go"}, map[string]int{"go": 1}},
		{[]string{"go", "is", "go"}, map[string]int{"go": 2, "is": 1}},
		{[]string{"a", "b", "a", "c", "a", "b"}, map[string]int{"a": 3, "b": 2, "c": 1}},
		{[]string{"Go", "go", "日本", "日本"}, map[string]int{"Go": 1, "go": 1, "日本": 2}},
	}
	for _, tt := range tests {
		got := CountWords(tt.words)
		if got == nil {
			t.Errorf("CountWords(%q) = nil, want a non-nil map", tt.words)
			continue
		}
		if len(got) != len(tt.want) {
			t.Errorf("CountWords(%q) has %d keys, want %d", tt.words, len(got), len(tt.want))
		}
		for w, n := range tt.want {
			c := got[w]
			if c == nil {
				t.Errorf("CountWords(%q)[%q] = nil, want a counter at %d", tt.words, w, n)
				continue
			}
			if c.Value() != n {
				t.Errorf("CountWords(%q)[%q].Value() = %d, want %d", tt.words, w, c.Value(), n)
			}
		}
	}
}
