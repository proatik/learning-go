# Receipt Writer

Write a receipt to **any** `io.Writer` with `fmt.Fprintf`, so the same code can print to the terminal or to a buffer in tests.

```go
type Item struct {
	Name  string
	Qty   int
	Cents int // price of ONE unit, in cents
}

func WriteReceipt(w io.Writer, items []Item) error
```

## Examples

`WriteReceipt(w, []Item{{"apple", 3, 50}, {"bread", 1, 249}})` writes:

```text
apple x3 = $1.50
bread x1 = $2.49
TOTAL = $3.99
```

`WriteReceipt(w, nil)` writes:

```text
(no items)
TOTAL = $0.00
```

## Rules
- One line per item: `<name> x<qty> = $<qty*cents as D.CC>`.
- Then a final line `TOTAL = $D.CC` with the sum of all line totals.
- For an empty (or nil) slice, write `(no items)` before the total.
- Every line ends with `\n`. Cents always have two digits (`$0.05`). You may assume no negative numbers.
- Return the **first** error returned by the writer (and stop writing); return `nil` on success.
