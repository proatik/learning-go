package crm

import (
	"reflect"
	"testing"
)

func TestNewCustomer(t *testing.T) {
	got := NewCustomer(7, "Ann Lee", "ann@example.com", "Oslo")
	want := Customer{ID: 7, Name: "Ann Lee", Email: "ann@example.com", Address: Address{City: "Oslo"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("NewCustomer(7, \"Ann Lee\", \"ann@example.com\", \"Oslo\") = %+v, want %+v", got, want)
	}
	if got.City != "Oslo" || got.Address.City != "Oslo" {
		t.Errorf("NewCustomer(...).City = %q, want \"Oslo\" (stored in the embedded Address)", got.City)
	}
}

func TestLabel(t *testing.T) {
	tests := []struct {
		c    Customer
		want string
	}{
		{NewCustomer(1, "Ann Lee", "ann@example.com", "Oslo"), "Ann Lee <ann@example.com> (Oslo)"},
		{Customer{Name: "Ann Lee", Address: Address{City: "Oslo"}}, "Ann Lee (Oslo)"},
		{Customer{Name: "Ann Lee", Email: "ann@example.com"}, "Ann Lee <ann@example.com>"},
		{Customer{Name: "Ann Lee", Address: Address{Street: "1 Main St", Zip: "0150"}}, "Ann Lee"},
		{Customer{Name: "Bo"}, "Bo"},
	}
	for _, tt := range tests {
		if got := Label(tt.c); got != tt.want {
			t.Errorf("Label(%+v) = %q, want %q", tt.c, got, tt.want)
		}
	}
}

func TestToJSONCustomer(t *testing.T) {
	full := Customer{
		ID:       7,
		Name:     "Ann Lee",
		Email:    "ann@example.com",
		Address:  Address{Street: "1 Main St", City: "Oslo", Zip: "0150"},
		Tags:     []string{"vip", "beta"},
		Password: "hunter2",
		notes:    "likes tea",
	}
	tests := []struct {
		name string
		c    Customer
		want string
	}{
		{"full", full, `{"id":7,"name":"Ann Lee","email":"ann@example.com","street":"1 Main St","city":"Oslo","zip":"0150","tags":["vip","beta"]}`},
		{"zero", Customer{}, `{"id":0,"name":"","street":"","city":""}`},
		{"empty tags", Customer{ID: 3, Name: "Cy", Tags: []string{}}, `{"id":3,"name":"Cy","street":"","city":""}`},
		{"password only", Customer{ID: 4, Name: "Di", Password: "secret"}, `{"id":4,"name":"Di","street":"","city":""}`},
		{"NewCustomer", NewCustomer(7, "Ann Lee", "ann@example.com", "Oslo"), `{"id":7,"name":"Ann Lee","email":"ann@example.com","street":"","city":"Oslo"}`},
	}
	for _, tt := range tests {
		got, err := ToJSON(tt.c)
		if err != nil {
			t.Errorf("%s: ToJSON(%+v) returned error %v", tt.name, tt.c, err)
			continue
		}
		if got != tt.want {
			t.Errorf("%s: ToJSON(%+v)\n got %s\nwant %s", tt.name, tt.c, got, tt.want)
		}
	}
}

func TestToJSONOrder(t *testing.T) {
	tests := []struct {
		o    Order
		want string
	}{
		{
			Order{ID: 1, Customer: Customer{ID: 2, Name: "Bo"}},
			`{"id":1,"customer":{"id":2,"name":"Bo","street":"","city":""}}`,
		},
		{
			Order{ID: 9, Customer: NewCustomer(5, "Ed", "", "Rome"), Note: "gift wrap"},
			`{"id":9,"customer":{"id":5,"name":"Ed","street":"","city":"Rome"},"note":"gift wrap"}`,
		},
		{Order{}, `{"id":0,"customer":{"id":0,"name":"","street":"","city":""}}`},
	}
	for _, tt := range tests {
		got, err := ToJSON(tt.o)
		if err != nil {
			t.Errorf("ToJSON(%+v) returned error %v", tt.o, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ToJSON(%+v)\n got %s\nwant %s", tt.o, got, tt.want)
		}
	}
}

func TestToJSONError(t *testing.T) {
	// A channel can't be encoded as JSON: the error must be passed on.
	got, err := ToJSON(make(chan int))
	if err == nil {
		t.Errorf("ToJSON(chan) = %q, nil; want an error", got)
	}
	if got != "" {
		t.Errorf("ToJSON(chan) = %q alongside the error, want \"\"", got)
	}
}
