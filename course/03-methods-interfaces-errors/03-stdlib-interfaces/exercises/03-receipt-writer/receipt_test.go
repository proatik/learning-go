package receipt

import (
	"bytes"
	"errors"
	"os"
	"strings"
	"testing"
)

func TestWriteReceipt(t *testing.T) {
	tests := []struct {
		items []Item
		want  string
	}{
		{
			[]Item{{"apple", 3, 50}, {"bread", 1, 249}},
			"apple x3 = $1.50\nbread x1 = $2.49\nTOTAL = $3.99\n",
		},
		{nil, "(no items)\nTOTAL = $0.00\n"},
		{[]Item{}, "(no items)\nTOTAL = $0.00\n"},
		{[]Item{{"gum", 1, 5}}, "gum x1 = $0.05\nTOTAL = $0.05\n"},
		{[]Item{{"free sample", 4, 0}}, "free sample x4 = $0.00\nTOTAL = $0.00\n"},
		{[]Item{{"tv", 2, 49999}, {"cable", 10, 199}}, "tv x2 = $999.98\ncable x10 = $19.90\nTOTAL = $1019.88\n"},
		{[]Item{{"none", 0, 300}}, "none x0 = $0.00\nTOTAL = $0.00\n"},
	}
	for _, tt := range tests {
		var buf bytes.Buffer
		if err := WriteReceipt(&buf, tt.items); err != nil {
			t.Errorf("WriteReceipt(%v) error = %v, want nil", tt.items, err)
		}
		if got := buf.String(); got != tt.want {
			t.Errorf("WriteReceipt(%v) = %q, want %q", tt.items, got, tt.want)
		}
	}
}

func TestWriteReceiptStringsBuilder(t *testing.T) {
	var sb strings.Builder
	items := []Item{{"tea", 2, 125}}
	want := "tea x2 = $2.50\nTOTAL = $2.50\n"
	if err := WriteReceipt(&sb, items); err != nil {
		t.Errorf("WriteReceipt(%v) error = %v, want nil", items, err)
	}
	if got := sb.String(); got != want {
		t.Errorf("WriteReceipt(%v) = %q, want %q", items, got, want)
	}
}

var errFull = errors.New("disk full")

// limitedWriter accepts `left` writes, then fails every write.
type limitedWriter struct {
	left  int
	calls int
}

func (w *limitedWriter) Write(p []byte) (int, error) {
	w.calls++
	if w.left == 0 {
		return 0, errFull
	}
	w.left--
	return len(p), nil
}

func TestWriteReceiptError(t *testing.T) {
	items := []Item{{"a", 1, 1}, {"b", 1, 1}, {"c", 1, 1}}
	for ok := 0; ok <= 3; ok++ {
		w := &limitedWriter{left: ok}
		err := WriteReceipt(w, items)
		if !errors.Is(err, errFull) {
			t.Errorf("WriteReceipt(writer failing after %d writes) error = %v, want %v", ok, err, errFull)
		}
		if w.calls != ok+1 {
			t.Errorf("WriteReceipt(writer failing after %d writes) made %d writes, want %d (stop at the first error)", ok, w.calls, ok+1)
		}
	}
}

func ExampleWriteReceipt() {
	WriteReceipt(os.Stdout, []Item{{"apple", 3, 50}, {"bread", 1, 249}})
	// Output:
	// apple x3 = $1.50
	// bread x1 = $2.49
	// TOTAL = $3.99
}
