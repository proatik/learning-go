package words

import (
	"math"
	"strings"
	"testing"
)

func TestNumberToWords(t *testing.T) {
	tests := []struct {
		n    int
		want string
	}{
		{0, "zero"},
		{1, "one"},
		{7, "seven"},
		{10, "ten"},
		{11, "eleven"},
		{12, "twelve"},
		{13, "thirteen"},
		{14, "fourteen"},
		{15, "fifteen"},
		{18, "eighteen"},
		{19, "nineteen"},
		{20, "twenty"},
		{21, "twenty-one"},
		{30, "thirty"},
		{40, "forty"},
		{42, "forty-two"},
		{50, "fifty"},
		{68, "sixty-eight"},
		{70, "seventy"},
		{80, "eighty"},
		{99, "ninety-nine"},
		{100, "one hundred"},
		{101, "one hundred one"},
		{110, "one hundred ten"},
		{115, "one hundred fifteen"},
		{123, "one hundred twenty-three"},
		{200, "two hundred"},
		{999, "nine hundred ninety-nine"},
		{1000, "one thousand"},
		{1001, "one thousand one"},
		{1010, "one thousand ten"},
		{1100, "one thousand one hundred"},
		{1234, "one thousand two hundred thirty-four"},
		{9000, "nine thousand"},
		{10_000, "ten thousand"},
		{12_345, "twelve thousand three hundred forty-five"},
		{20_020, "twenty thousand twenty"},
		{80_008, "eighty thousand eight"},
		{100_000, "one hundred thousand"},
		{100_001, "one hundred thousand one"},
		{110_011, "one hundred ten thousand eleven"},
		{505_050, "five hundred five thousand fifty"},
		{999_999, "nine hundred ninety-nine thousand nine hundred ninety-nine"},
	}
	for _, tt := range tests {
		got, err := NumberToWords(tt.n)
		if err != nil || got != tt.want {
			t.Errorf("NumberToWords(%d) = %q, %v, want %q, nil", tt.n, got, err, tt.want)
		}
	}
}

func TestNumberToWordsFormat(t *testing.T) {
	for n := range 1_000_000 {
		got, err := NumberToWords(n)
		if err != nil || got == "" {
			t.Fatalf("NumberToWords(%d) = %q, %v, want words, nil", n, got, err)
		}
		if n != 0 && strings.Contains(got, "zero") {
			t.Fatalf("NumberToWords(%d) = %q, must not contain \"zero\"", n, got)
		}
		if strings.Contains(got, "  ") || strings.TrimSpace(got) != got || strings.Contains(got, " and ") {
			t.Fatalf("NumberToWords(%d) = %q, has bad spacing or \"and\"", n, got)
		}
	}
}

func TestNumberToWordsInvalid(t *testing.T) {
	for _, n := range []int{-1, -999, 1_000_000, 1_000_001, math.MaxInt, math.MinInt} {
		if got, err := NumberToWords(n); err == nil || got != "" {
			t.Errorf("NumberToWords(%d) = %q, %v, want \"\", error", n, got, err)
		}
	}
}
