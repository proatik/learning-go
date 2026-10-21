package people

import (
	"slices"
	"testing"
)

func sample() []Person {
	return []Person{{"Ana", 30}, {"Ben", 41}, {"Cy", 41}, {"Di", 7}}
}

func TestBirthday(t *testing.T) {
	p := Person{"Ana", 30}
	Birthday(&p)
	if p.Age != 31 {
		t.Errorf("Birthday(&{Ana 30}): Age = %d, want 31", p.Age)
	}
	Birthday(&p)
	if p.Age != 32 {
		t.Errorf("Birthday twice from 30: Age = %d, want 32", p.Age)
	}
	if p.Name != "Ana" {
		t.Errorf("Birthday changed Name to %q", p.Name)
	}
}

func TestBirthdayNil(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Birthday(nil) panicked: %v", r)
		}
	}()
	Birthday(nil)
}

func TestBirthdayAll(t *testing.T) {
	people := sample()
	BirthdayAll(people)
	want := []Person{{"Ana", 31}, {"Ben", 42}, {"Cy", 42}, {"Di", 8}}
	if !slices.Equal(people, want) {
		t.Errorf("BirthdayAll: people = %v, want %v", people, want)
	}
	BirthdayAll(nil) // must not panic
	BirthdayAll([]Person{})
}

func TestFind(t *testing.T) {
	people := sample()
	for i, want := range people {
		p := Find(people, want.Name)
		if p == nil {
			t.Fatalf("Find(people, %q) = nil, want &people[%d]", want.Name, i)
		}
		if p != &people[i] {
			t.Errorf("Find(people, %q) = %p (%v), want &people[%d] = %p", want.Name, p, *p, i, &people[i])
		}
	}
}

func TestFindModifiesSlice(t *testing.T) {
	people := sample()
	p := Find(people, "Ben")
	if p == nil {
		t.Fatalf("Find(people, %q) = nil", "Ben")
	}
	p.Age = 99
	if people[1].Age != 99 {
		t.Errorf("after Find(people, \"Ben\").Age = 99: people[1] = %v, want {Ben 99} (return &people[i], not a copy)", people[1])
	}
}

func TestFindFirstDuplicate(t *testing.T) {
	people := []Person{{"Ana", 1}, {"Ana", 2}}
	if p := Find(people, "Ana"); p != &people[0] {
		t.Errorf("Find with duplicate names did not return the first element")
	}
}

func TestFindMissing(t *testing.T) {
	for _, name := range []string{"Zed", "", "ana"} {
		if p := Find(sample(), name); p != nil {
			t.Errorf("Find(people, %q) = %v, want nil", name, *p)
		}
	}
	if p := Find(nil, "Ana"); p != nil {
		t.Errorf("Find(nil, \"Ana\") = %v, want nil", *p)
	}
}

func TestOldest(t *testing.T) {
	tests := []struct {
		people []Person
		index  int
	}{
		{sample(), 1}, // Ben and Cy tie at 41: first one wins
		{[]Person{{"Solo", 0}}, 0},
		{[]Person{{"A", 5}, {"B", 6}, {"C", 99}}, 2},
		{[]Person{{"A", 3}, {"B", 2}, {"C", 1}}, 0},
	}
	for _, tt := range tests {
		got := Oldest(tt.people)
		if got == nil {
			t.Fatalf("Oldest(%v) = nil, want &people[%d]", tt.people, tt.index)
		}
		if got != &tt.people[tt.index] {
			t.Errorf("Oldest(%v) = %v, want &people[%d] (%v)", tt.people, *got, tt.index, tt.people[tt.index])
		}
	}
}

func TestOldestEmpty(t *testing.T) {
	if got := Oldest(nil); got != nil {
		t.Errorf("Oldest(nil) = %v, want nil", *got)
	}
	if got := Oldest([]Person{}); got != nil {
		t.Errorf("Oldest([]) = %v, want nil", *got)
	}
}

func TestOldestThenBirthday(t *testing.T) {
	people := sample()
	Birthday(Oldest(people))
	if people[1].Age != 42 || people[2].Age != 41 {
		t.Errorf("Birthday(Oldest(people)): people = %v, want Ben 42 and Cy 41", people)
	}
}
