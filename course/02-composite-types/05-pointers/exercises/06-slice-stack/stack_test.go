package stack

import (
	"slices"
	"testing"
)

func TestPush(t *testing.T) {
	var s []int
	Push(&s, 1)
	Push(&s, 2)
	Push(&s, 3)
	if want := []int{1, 2, 3}; !slices.Equal(s, want) {
		t.Errorf("Push 1, 2, 3 onto nil: s = %v, want %v", s, want)
	}
}

func TestPushAll(t *testing.T) {
	s := []int{9}
	PushAll(&s, 1, 2, 3)
	if want := []int{9, 1, 2, 3}; !slices.Equal(s, want) {
		t.Errorf("PushAll(&[9], 1, 2, 3): s = %v, want %v", s, want)
	}
	PushAll(&s)
	if want := []int{9, 1, 2, 3}; !slices.Equal(s, want) {
		t.Errorf("PushAll with no values changed s to %v", s)
	}
}

func TestPop(t *testing.T) {
	s := []int{1, 2, 3}
	for _, want := range []int{3, 2, 1} {
		got, ok := Pop(&s)
		if !ok || got != want {
			t.Fatalf("Pop = (%d, %v), want (%d, true); s is now %v", got, ok, want, s)
		}
	}
	if len(s) != 0 {
		t.Errorf("after popping everything: len(s) = %d, want 0", len(s))
	}
}

func TestPopEmpty(t *testing.T) {
	var s []int
	if got, ok := Pop(&s); ok || got != 0 {
		t.Errorf("Pop(&nil) = (%d, %v), want (0, false)", got, ok)
	}
	s = []int{}
	if got, ok := Pop(&s); ok || got != 0 {
		t.Errorf("Pop(&[]) = (%d, %v), want (0, false)", got, ok)
	}
	if len(s) != 0 {
		t.Errorf("Pop on empty changed len to %d", len(s))
	}
}

func TestPopShrinksCallerSlice(t *testing.T) {
	s := []int{4, 5}
	Pop(&s)
	if want := []int{4}; !slices.Equal(s, want) {
		t.Errorf("Pop(&[4 5]): s = %v, want %v (assign through the pointer: *s = (*s)[:n-1])", s, want)
	}
}

func TestPeek(t *testing.T) {
	s := []int{1, 2, 3}
	got, ok := Peek(s)
	if !ok || got != 3 {
		t.Errorf("Peek([1 2 3]) = (%d, %v), want (3, true)", got, ok)
	}
	if want := []int{1, 2, 3}; !slices.Equal(s, want) {
		t.Errorf("Peek changed s to %v", s)
	}
	if got, ok := Peek(nil); ok || got != 0 {
		t.Errorf("Peek(nil) = (%d, %v), want (0, false)", got, ok)
	}
}

func TestMixed(t *testing.T) {
	var s []int
	Push(&s, 10)
	PushAll(&s, 20, 30)
	if v, ok := Pop(&s); !ok || v != 30 {
		t.Fatalf("Pop = (%d, %v), want (30, true)", v, ok)
	}
	Push(&s, 40)
	if v, ok := Peek(s); !ok || v != 40 {
		t.Errorf("Peek = (%d, %v), want (40, true)", v, ok)
	}
	if want := []int{10, 20, 40}; !slices.Equal(s, want) {
		t.Errorf("stack = %v, want %v", s, want)
	}
	var popped []int
	for {
		v, ok := Pop(&s)
		if !ok {
			break
		}
		popped = append(popped, v)
	}
	if want := []int{40, 20, 10}; !slices.Equal(popped, want) {
		t.Errorf("popped in order %v, want %v", popped, want)
	}
}

func TestManyPushes(t *testing.T) {
	var s []int
	for i := range 1000 {
		Push(&s, i)
	}
	if len(s) != 1000 || s[0] != 0 || s[999] != 999 {
		t.Fatalf("after 1000 pushes: len = %d, want 1000 with s[0]=0 and s[999]=999", len(s))
	}
	for i := 999; i >= 0; i-- {
		if v, ok := Pop(&s); !ok || v != i {
			t.Fatalf("Pop = (%d, %v), want (%d, true)", v, ok, i)
		}
	}
}
