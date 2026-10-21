# Customer Records: Embedding & JSON Tags

`Customer` **embeds** `Address`, so `c.City` works directly. Your job: write three functions, and add **struct tags** so `json.Marshal` produces exactly the JSON below.

```go
type Address struct {
	Street string
	City   string
	Zip    string
}

type Customer struct {
	ID       int
	Name     string
	Email    string
	Address  // embedded
	Tags     []string
	Password string
	notes    string
}

type Order struct {
	ID       int
	Customer Customer // a named field: nested, not embedded
	Note     string
}

func NewCustomer(id int, name, email, city string) Customer
func Label(c Customer) string
func ToJSON(v any) (string, error) // json.Marshal(v) as a string
```

## Required JSON

| Field | JSON key | Rule |
|---|---|---|
| `Customer.ID`, `Customer.Name` | `"id"`, `"name"` | always present |
| `Customer.Email` | `"email"` | omitted when empty |
| `Address.Street`, `Address.City` | `"street"`, `"city"` | always present, **flattened** into the customer object |
| `Address.Zip` | `"zip"` | omitted when empty |
| `Customer.Tags` | `"tags"` | omitted when empty or nil |
| `Customer.Password` | — | **never** encoded |
| `Customer.notes` | — | never encoded (unexported, no tag needed) |
| `Order.ID`, `Order.Customer` | `"id"`, `"customer"` | always; the customer is a nested object |
| `Order.Note` | `"note"` | omitted when empty |

```go
c := NewCustomer(7, "Ann Lee", "ann@example.com", "Oslo")
c.Tags = []string{"vip"}
c.Password = "hunter2"
ToJSON(c)
// {"id":7,"name":"Ann Lee","email":"ann@example.com","street":"","city":"Oslo","tags":["vip"]}

ToJSON(Order{ID: 1, Customer: Customer{ID: 2, Name: "Bo"}})
// {"id":1,"customer":{"id":2,"name":"Bo","street":"","city":""}}
```

## Functions

| Call | Result |
|---|---|
| `NewCustomer(7, "Ann Lee", "ann@example.com", "Oslo")` | `Customer{ID: 7, Name: "Ann Lee", Email: "ann@example.com", City: "Oslo"}` |
| `Label(c)` | `"Ann Lee <ann@example.com> (Oslo)"` |
| `Label` with no email | `"Ann Lee (Oslo)"` |
| `Label` with no city | `"Ann Lee <ann@example.com>"` |
| `Label` with neither | `"Ann Lee"` |

## Rules
- Do not rename or reorder fields: JSON keys follow field order.
- In `NewCustomer`, try the Go 1.27 promoted-field keys: `Customer{..., City: city}`.
- `ToJSON` returns `json.Marshal`'s error unchanged (and `""` with it).
