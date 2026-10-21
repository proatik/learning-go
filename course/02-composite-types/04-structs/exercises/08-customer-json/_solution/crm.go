package crm

import "encoding/json"

// Address is a postal address.
type Address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	Zip    string `json:"zip,omitempty"`
}

// Customer is a customer record. Address is embedded, so its fields are
// promoted: c.City is c.Address.City.
type Customer struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email,omitempty"`
	Address           // embedded without a tag: its fields are flattened into the object
	Tags     []string `json:"tags,omitempty"`
	Password string   `json:"-"`
	notes    string   // unexported: json ignores it
}

// Order is an order placed by a customer.
type Order struct {
	ID       int      `json:"id"`
	Customer Customer `json:"customer"` // a named field: encoded as a nested object
	Note     string   `json:"note,omitempty"`
}

// NewCustomer returns a customer with the given ID, name, email and city.
func NewCustomer(id int, name, email, city string) Customer {
	// Go 1.27: the promoted field City can be used as a key directly.
	return Customer{ID: id, Name: name, Email: email, City: city}
}

// Label returns "Name <email> (City)", leaving out the parts that are empty.
func Label(c Customer) string {
	s := c.Name
	if c.Email != "" {
		s += " <" + c.Email + ">"
	}
	if c.City != "" { // promoted from Address
		s += " (" + c.City + ")"
	}
	return s
}

// ToJSON returns the JSON encoding of v.
func ToJSON(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
