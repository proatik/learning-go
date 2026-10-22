package rle

import "testing"

// FuzzRoundTrip checks that Decode(Encode(s)) == s.
func FuzzRoundTrip(f *testing.F) {
	// TODO: add seeds with f.Add that reach every tricky case.
	f.Fuzz(func(t *testing.T, s string) {
		// TODO: skip invalid UTF-8, then check the round trip.
		_ = Encode(s)
	})
}

// TODO: FuzzDecode — Decode must never panic on any input.
// TODO: TestDecodeErrors — a table of malformed inputs.
// TODO: ExampleEncode with an // Output: comment.
// TODO: BenchmarkEncode using b.Loop.
