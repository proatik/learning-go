package size

// TODO: use iota (and _ to skip the first value) so that
// KB = 1 << 10, MB = 1 << 20, GB = 1 << 30, TB = 1 << 40.
const (
	KB = 0
	MB = 0
	GB = 0
	TB = 0
)

// HumanSize formats a byte count using the largest fitting unit
// (B, KB, MB, GB or TB), with one decimal for units above bytes.
func HumanSize(n int64) string {
	return "" // TODO
}
