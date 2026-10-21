package crm

// TODO: add json struct tags to Address, Customer and Order (see README).

// Address is a postal address.
type Address struct {
	Street string
	City   string
	Zip    string
}

// Customer is a customer record. Address is embedded, so its fields are
// promoted: c.City is c.Address.City.
type Customer struct {
	ID    int
	Name  string
	Email string
	Address
	Tags     []string
	Password string
	notes    string
}

// Order is an order placed by a customer.
type Order struct {
	ID       int
	Customer Customer
	Note     string
}

// NewCustomer returns a customer with the given ID, name, email and city.
func NewCustomer(id int, name, email, city string) Customer {
	// TODO
	return Customer{}
}

// Label returns "Name <email> (City)", leaving out the parts that are empty.
func Label(c Customer) string {
	// TODO
	return ""
}

// ToJSON returns the JSON encoding of v.
func ToJSON(v any) (string, error) {
	// TODO
	return "", nil
}
