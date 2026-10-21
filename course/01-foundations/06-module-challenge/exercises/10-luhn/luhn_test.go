package luhn

import (
	"strconv"
	"testing"
)

func TestValid(t *testing.T) {
	tests := []struct {
		number string
		want   bool
	}{
		{"4539 3195 0343 6467", true},
		{"4539319503436467", true},
		{"8273 1232 7352 0569", false},
		{"79927398713", true},
		{"79927398710", false},
		{"059", true},
		{"59", true},
		{"0 0", true},
		{" 0 ", false}, // only one digit
		{"0", false},
		{"", false},
		{"   ", false},
		{"091", true},
		{"9999999999 9999999999 9999999999 9999999999", true},
		{"055 444 285", true},
		{"055-444-285", false},
		{"055a 444 285", false},
		{"055 444 28٥", false}, // Arabic-Indic digit is not an ASCII digit
		{"1234567812345670", true},
		{"1234567812345678", false},
	}
	for _, tt := range tests {
		if got := Valid(tt.number); got != tt.want {
			t.Errorf("Valid(%q) = %v, want %v", tt.number, got, tt.want)
		}
	}
}

func TestCheckDigit(t *testing.T) {
	tests := []struct {
		payload string
		want    int
	}{
		{"7992739871", 3},
		{"453931950343646", 7},
		{"123456781234567", 0},
		{"05", 9},
		{"0", 0},
		{"1", 8},
		{"4 5 3 9", 3},
	}
	for _, tt := range tests {
		got, err := CheckDigit(tt.payload)
		if err != nil {
			t.Errorf("CheckDigit(%q) returned error %v, want %d", tt.payload, err, tt.want)
			continue
		}
		if got != tt.want {
			t.Errorf("CheckDigit(%q) = %d, want %d", tt.payload, got, tt.want)
		}
	}
}

func TestCheckDigitErrors(t *testing.T) {
	for _, payload := range []string{"", "  ", "12a4", "12-4"} {
		if _, err := CheckDigit(payload); err == nil {
			t.Errorf("CheckDigit(%q) returned nil error, want an error", payload)
		}
	}
}

// TestCheckDigitMakesValid checks the two functions agree for many payloads.
func TestCheckDigitMakesValid(t *testing.T) {
	for n := range 2000 {
		payload := strconv.Itoa(n * 7919)
		d, err := CheckDigit(payload)
		if err != nil {
			t.Fatalf("CheckDigit(%q) returned error %v", payload, err)
		}
		if d < 0 || d > 9 {
			t.Fatalf("CheckDigit(%q) = %d, want a single digit", payload, d)
		}
		full := payload + strconv.Itoa(d)
		if !Valid(full) {
			t.Fatalf("Valid(%q) = false, but CheckDigit(%q) said %d", full, payload, d)
		}
	}
}
