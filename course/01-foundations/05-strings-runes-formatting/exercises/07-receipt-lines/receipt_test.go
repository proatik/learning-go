package receipt

import (
	"fmt"
	"testing"
	"unicode/utf8"
)

func TestHeader(t *testing.T) {
	if got, want := Header(), "Item         Qty     Price      Total"; got != want {
		t.Errorf("Header() = %q, want %q", got, want)
	}
}

func TestSeparator(t *testing.T) {
	if got, want := Separator(), "-------------------------------------"; got != want {
		t.Errorf("Separator() = %q, want %q", got, want)
	}
}

func TestLine(t *testing.T) {
	tests := []struct {
		item string
		qty  int
		unit float64
		want string
	}{
		{"Coffee", 2, 3.5, "Coffee         2 x    3.50 =     7.00"},
		{"Café", 1, 4.25, "Café           1 x    4.25 =     4.25"},
		{"Croissant au beurre", 3, 2.1, "Croissant au   3 x    2.10 =     6.30"},
		{"抹茶ラテ", 10, 5, "抹茶ラテ          10 x    5.00 =    50.00"},
		{"Water", 0, 1.99, "Water          0 x    1.99 =     0.00"},
		{"Big order", 120, 99.99, "Big order    120 x   99.99 = 11998.80"},
		{"", 1, 0.5, "               1 x    0.50 =     0.50"},
	}
	for _, tt := range tests {
		if got := Line(tt.item, tt.qty, tt.unit); got != tt.want {
			t.Errorf("Line(%q, %d, %v) =\n  %q\nwant\n  %q", tt.item, tt.qty, tt.unit, got, tt.want)
		}
	}
}

func TestLineWidth(t *testing.T) {
	for _, item := range []string{"Tea", "Café crème brûlée", "世界世界世界世界世界世界世界"} {
		if n := utf8.RuneCountInString(Line(item, 1, 1)); n != 37 {
			t.Errorf("Line(%q, 1, 1) is %d runes wide, want 37", item, n)
		}
	}
}

func TestTotal(t *testing.T) {
	tests := []struct {
		total float64
		want  string
	}{
		{17.55, "TOTAL                           17.55"},
		{0, "TOTAL                            0.00"},
		{1234.5, "TOTAL                         1234.50"},
	}
	for _, tt := range tests {
		if got := Total(tt.total); got != tt.want {
			t.Errorf("Total(%v) = %q, want %q", tt.total, got, tt.want)
		}
	}
}

func Example() {
	fmt.Println(Header())
	fmt.Println(Separator())
	fmt.Println(Line("Coffee", 2, 3.5))
	fmt.Println(Line("Café", 1, 4.25))
	fmt.Println(Line("Croissant au beurre", 3, 2.1))
	fmt.Println(Total(17.55))
	// Output:
	// Item         Qty     Price      Total
	// -------------------------------------
	// Coffee         2 x    3.50 =     7.00
	// Café           1 x    4.25 =     4.25
	// Croissant au   3 x    2.10 =     6.30
	// TOTAL                           17.55
}
