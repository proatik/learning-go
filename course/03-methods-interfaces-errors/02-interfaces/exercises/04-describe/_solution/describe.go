package kinds

import (
	"fmt"
	"unicode/utf8"
)

// Describe returns a short description of v based on its dynamic type.
func Describe(v any) string {
	switch x := v.(type) {
	case nil:
		return "nil"
	case int:
		return fmt.Sprintf("int %d", x)
	case int8, int16, int32, int64: // x is still an any here
		return fmt.Sprintf("sized int %v", x)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("unsigned %v", x)
	case float64:
		return fmt.Sprintf("float64 %v", x)
	case string:
		return fmt.Sprintf("string %q with %d runes", x, utf8.RuneCountInString(x))
	case bool:
		return fmt.Sprintf("bool %t", x)
	case []int:
		return fmt.Sprintf("[]int with %d elements", len(x))
	case map[string]int:
		return fmt.Sprintf("map[string]int with %d keys", len(x))
	case fmt.Stringer:
		return "stringer " + x.String()
	default:
		return fmt.Sprintf("unknown %T", x)
	}
}
