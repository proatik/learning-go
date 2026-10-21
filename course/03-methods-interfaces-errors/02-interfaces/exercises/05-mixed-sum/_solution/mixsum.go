package mixsum

// Valuer is any type that knows its numeric value.
type Valuer interface {
	Value() float64
}

// Sum adds every numeric element of vals (int, int64, float32, float64,
// non-nil *int and Valuer) and reports how many elements were skipped.
func Sum(vals []any) (total float64, skipped int) {
	for _, v := range vals {
		switch x := v.(type) {
		case int:
			total += float64(x)
		case int64:
			total += float64(x)
		case float32:
			total += float64(x)
		case float64:
			total += x
		case *int:
			if x == nil {
				skipped++
				continue
			}
			total += float64(*x)
		case Valuer:
			total += x.Value()
		default:
			skipped++
		}
	}
	return total, skipped
}

// Strings returns the string elements of vals in order.
func Strings(vals []any) []string {
	var out []string
	for _, v := range vals {
		if s, ok := v.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
