package stack

import (
	"slices"
	"testing"
)

func TestZeroValue(t *testing.T) {
	var s Stack
	if got := s.Len(); got != 0 {
		t.Errorf("empty Stack Len() = %d, want 0", got)
	}
	if v, ok := s.Pop(); v != 0 || ok {
		t.Errorf("empty Stack Pop() = (%d, %v), want (0, false)", v, ok)
	}
	if v, ok := s.Peek(); v != 0 || ok {
		t.Errorf("empty Stack Peek() = (%d, %v), want (0, false)", v, ok)
	}
	if got := s.Values(); len(got) != 0 {
		t.Errorf("empty Stack Values() = %v, want []", got)
	}
}

func TestPushValues(t *testing.T) {
	tests := []struct {
		push []int
		want []int
	}{
		{[]int{1}, []int{1}},
		{[]int{1, 2}, []int{1, 2}},
		{[]int{5, -3, 0, 9}, []int{5, -3, 0, 9}},
	}
	for _, tt := range tests {
		var s Stack
		for _, v := range tt.push {
			s.Push(v)
		}
		if got := s.Values(); !slices.Equal(got, tt.want) {
			t.Errorf("Push %v: Values() = %v, want %v", tt.push, got, tt.want)
		}
		if got := s.Len(); got != len(tt.want) {
			t.Errorf("Push %v: Len() = %d, want %d", tt.push, got, len(tt.want))
		}
	}
}

func TestPeekPop(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push(2)
	s.Push(3)
	if v, ok := s.Peek(); v != 3 || !ok {
		t.Errorf("Peek() on [1 2 3] = (%d, %v), want (3, true)", v, ok)
	}
	if got := s.Len(); got != 3 {
		t.Errorf("Peek() changed the stack: Len() = %d, want 3", got)
	}
	for _, want := range []int{3, 2, 1} {
		if v, ok := s.Pop(); v != want || !ok {
			t.Errorf("Pop() = (%d, %v), want (%d, true)", v, ok, want)
		}
	}
	if v, ok := s.Pop(); v != 0 || ok {
		t.Errorf("Pop() on emptied stack = (%d, %v), want (0, false)", v, ok)
	}
	s.Push(7)
	if v, ok := s.Pop(); v != 7 || !ok {
		t.Errorf("Push(7) after emptying, Pop() = (%d, %v), want (7, true)", v, ok)
	}
}

func TestValuesIsCopy(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push(2)
	vals := s.Values()
	if len(vals) != 2 {
		t.Fatalf("Values() = %v, want [1 2]", vals)
	}
	vals[0] = 99
	if got := s.Values(); !slices.Equal(got, []int{1, 2}) {
		t.Errorf("modifying the Values() result changed the stack: Values() = %v, want [1 2]", got)
	}
}

func TestEvalRPN(t *testing.T) {
	tests := []struct {
		tokens []string
		want   int
		ok     bool
	}{
		{[]string{"2", "3", "+", "4", "*"}, 20, true},
		{[]string{"10", "3", "-"}, 7, true},
		{[]string{"3", "10", "-"}, -7, true},
		{[]string{"7", "2", "/"}, 3, true},
		{[]string{"-7", "2", "/"}, -3, true},
		{[]string{"42"}, 42, true},
		{[]string{"-3", "-4", "*"}, 12, true},
		{[]string{"5", "1", "2", "+", "4", "*", "+", "3", "-"}, 14, true},
		{[]string{"1", "+"}, 0, false},
		{[]string{"+"}, 0, false},
		{[]string{"4", "0", "/"}, 0, false},
		{[]string{"1", "2"}, 0, false},
		{[]string{"1", "x", "+"}, 0, false},
		{[]string{"2", "3", "%"}, 0, false},
		{nil, 0, false},
		{[]string{}, 0, false},
	}
	for _, tt := range tests {
		got, ok := EvalRPN(tt.tokens)
		if got != tt.want || ok != tt.ok {
			t.Errorf("EvalRPN(%q) = (%d, %v), want (%d, %v)", tt.tokens, got, ok, tt.want, tt.ok)
		}
	}
}
