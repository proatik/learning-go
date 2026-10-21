package acronym

import "testing"

func TestAcronym(t *testing.T) {
	tests := []struct {
		phrase string
		want   string
	}{
		{"Portable Network Graphics", "PNG"},
		{"as soon as possible", "ASAP"},
		{"Complementary metal-oxide semiconductor", "CMOS"},
		{"  Ruby   on\tRails ", "ROR"},
		{"Rock & Roll", "RR"},
		{"The Road _Not_ Taken", "TRNT"},
		{"élan vital", "ÉV"},
		{"über straße", "ÜS"},
		{"Go", "G"},
		{"", ""},
		{"   ", ""},
		{"- - -", ""},
		{"HyperText Markup Language", "HML"},
		{"2024 summer olympics", "SO"},
	}
	for _, tt := range tests {
		if got := Acronym(tt.phrase); got != tt.want {
			t.Errorf("Acronym(%q) = %q, want %q", tt.phrase, got, tt.want)
		}
	}
}
