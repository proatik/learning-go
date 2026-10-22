# Bug Hunt: Checkout Discounts

A checkout `Total` with volume discounts and promo codes is already written. Turn its rules into **tests** in `checkout_test.go` that pass on the real code and catch all **12** planted bugs — most of them hide at a boundary.

```go
// Provided in check_test.go (read it in the Tests tab).
var (
	ErrEmptyCart      = errors.New("checkout: empty cart")
	ErrBadLine        = errors.New("checkout: bad line")
	ErrUnknownCode    = errors.New("checkout: unknown code")
	ErrCodeNotAllowed = errors.New("checkout: code not allowed")
)

type Line struct {
	SKU   string
	Cents int64
	Qty   int
}

func Total(lines []Line, code string) (int64, error)
```

## The rules you are testing

| # | Rule |
|---|---|
| 1 | An empty (or `nil`) cart → `ErrEmptyCart`. |
| 2 | A line with `Qty < 1` or `Cents < 0` → `ErrBadLine`. Free items (`Cents == 0`) are fine. |
| 3 | **Bulk:** a line with `Qty >= 10` costs 10% less: `line * 90 / 100`, rounded **down**. Only that line is discounted. |
| 4 | Code `""`: no promo. Codes are **case-insensitive** (`"save5"` = `"SAVE5"`). |
| 5 | `SAVE5`: 500 off, allowed only if the total **after** bulk discounts is **at least** 2500; otherwise `ErrCodeNotAllowed`. |
| 6 | `HALF`: takes half (rounded down) of **one unit** of the **most expensive** item off the total. |
| 7 | Any other code → `ErrUnknownCode`. |

## Examples

| Code | Result |
|---|---|
| `Total([{a 250 2} {b 100 1}], "")` | `600, nil` |
| `Total([{a 30 10}], "")` | `270, nil` — bulk |
| `Total([{a 2500 1}], "SAVE5")` | `2000, nil` |
| `Total([{a 300 2} {b 1000 3}], "HALF")` | `3100, nil` — 3600 − 500 |
| `Total(nil, "")` | `0`, `ErrEmptyCart` |

## Rules
- Write at least one `Test…` function. Using **subtests** and a `t.Helper()` function for the repeated checks is strongly recommended.
- Check errors with `errors.Is(err, ErrX)`, never by comparing text.
- Each bug changes exactly one rule above. Read the bug names in `check_test.go`, then design an input where the buggy and the correct answer **differ** — boundaries (`9`/`10`/`11` units, `2499`/`2500`), rounding, carts with **two** lines, mixed-case codes.
- A failing run lists each surviving bug as `Bug "…" = not caught`.
