package chunk

import (
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		in   []int
		size int
		want [][]int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{[]int{1, 2, 3}, 5, [][]int{{1, 2, 3}}},
		{[]int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{[]int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
		{[]int{1, 2, 3, 4, 5, 6, 7}, 3, [][]int{{1, 2, 3}, {4, 5, 6}, {7}}},
		{[]int{42}, 2, [][]int{{42}}},
	}
	for _, tt := range tests {
		got := Chunk(tt.in, tt.size)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Chunk(%v, %d) = %v, want %v", tt.in, tt.size, got, tt.want)
		}
	}
}

func TestChunkEmpty(t *testing.T) {
	for _, in := range [][]int{nil, {}} {
		if got := Chunk(in, 3); len(got) != 0 {
			t.Errorf("Chunk(%v, 3) = %v, want no chunks", in, got)
		}
	}
}

func TestChunkInvalidSize(t *testing.T) {
	for _, size := range []int{0, -1, -5} {
		if got := Chunk([]int{1, 2}, size); got != nil {
			t.Errorf("Chunk([1 2], %d) = %v, want nil", size, got)
		}
	}
}

func TestChunkSharesBackingArray(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	chunks := Chunk(s, 2)
	if len(chunks) != 3 || len(chunks[1]) == 0 {
		t.Fatalf("Chunk([1 2 3 4 5], 2) = %v, want [[1 2] [3 4] [5]]", chunks)
	}
	chunks[1][0] = 30
	if s[2] != 30 {
		t.Errorf("chunks must be sub-slices of s (no copying): writing chunks[1][0] did not change s[2]; s = %v", s)
	}
}

func TestChunkCapEqualsLen(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	chunks := Chunk(s, 2)
	if len(chunks) != 3 {
		t.Fatalf("Chunk([1 2 3 4 5], 2) = %v, want 3 chunks", chunks)
	}
	for i, c := range chunks {
		if cap(c) != len(c) {
			t.Errorf("chunk %d = %v has len %d but cap %d; want cap == len (use s[lo:hi:hi])", i, c, len(c), cap(c))
		}
	}
	_ = append(chunks[0], 99) // must not overwrite s[2]
	if s[2] != 3 || chunks[1][0] != 3 {
		t.Errorf("appending to chunks[0] overwrote the next chunk: s = %v, chunks[1] = %v", s, chunks[1])
	}
}
