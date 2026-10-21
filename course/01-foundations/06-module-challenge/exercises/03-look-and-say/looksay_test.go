package looksay

import "testing"

func TestNext(t *testing.T) {
	tests := []struct {
		s, want string
	}{
		{"1", "11"},
		{"11", "21"},
		{"21", "1211"},
		{"1211", "111221"},
		{"111221", "312211"},
		{"0", "10"},
		{"9", "19"},
		{"1112", "3112"},
		{"3333333333", "103"},
		{"222222222222", "122"},
		{"1223334444", "11223344"},
		{"90009", "193019"},
	}
	for _, tt := range tests {
		got, err := Next(tt.s)
		if err != nil || got != tt.want {
			t.Errorf("Next(%q) = %q, %v, want %q, nil", tt.s, got, err, tt.want)
		}
	}

	for _, s := range []string{"", "a", "12a", "1 1", "-1", "１", "1.0"} {
		if got, err := Next(s); err == nil {
			t.Errorf("Next(%q) = %q, nil, want an error", s, got)
		}
	}
}

func TestTerm(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{1, "1"},
		{2, "11"},
		{3, "21"},
		{4, "1211"},
		{5, "111221"},
		{6, "312211"},
		{7, "13112221"},
		{10, "13211311123113112211"},
	}
	for _, tt := range tests {
		got, err := Term(tt.n)
		if err != nil || got != tt.want {
			t.Errorf("Term(%d) = %q, %v, want %q, nil", tt.n, got, err, tt.want)
		}
	}

	if got, err := Term(20); err != nil || len(got) != 302 {
		t.Errorf("Term(20) has length %d (err %v), want 302", len(got), err)
	}

	for _, n := range []int{0, -1, -100} {
		if got, err := Term(n); err == nil {
			t.Errorf("Term(%d) = %q, nil, want an error", n, got)
		}
	}
}
