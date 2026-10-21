# Key-Value Store

Build an in-memory store out of small, embedded interfaces, then write helpers that accept only the behaviour they need.

```go
type Getter interface{ Get(key string) (string, bool) }
type Setter interface{ Set(key, value string) }
type Deleter interface{ Delete(key string) bool }

type Store interface {
	Getter
	Setter
	Deleter
	Len() int
}

type MemStore struct{ /* unexported fields */ }

func NewMemStore() *MemStore
func (m *MemStore) Get(key string) (string, bool)
func (m *MemStore) Set(key, value string)
func (m *MemStore) Delete(key string) bool
func (m *MemStore) Len() int
func (m *MemStore) Keys() []string

func GetOr(g Getter, key, def string) string
func CopyKeys(dst Setter, src Getter, keys []string) int
func ReadOnly(g Getter) Getter
```

## Examples

| Code | Result |
|---|---|
| `m := NewMemStore(); m.Set("a", "1"); m.Get("a")` | `"1", true` |
| `m.Get("zzz")` | `"", false` |
| `m.Delete("a")` / again | `true` / `false` |
| `m.Set("b", "2"); m.Set("a", "1"); m.Keys()` | `["a" "b"]` |
| `GetOr(m, "missing", "n/a")` | `"n/a"` |
| `CopyKeys(dst, src, []string{"a", "x"})` (only `"a"` in `src`) | `1` |
| `ro := ReadOnly(m); ro.Get("a")` | same as `m.Get("a")` |
| `_, ok := ro.(Setter)` | `ok == false` |

## Rules
- `*MemStore` must satisfy `Store` (add `var _ Store = (*MemStore)(nil)`).
- The **zero value** `MemStore{}` must be usable too: `Set` creates the map on first use; `Get`, `Delete`, `Len`, `Keys` work on an empty store.
- `Set` overwrites an existing key. `Delete` reports whether the key existed.
- `Keys` returns the keys sorted ascending (empty store → empty, `nil` is fine).
- `GetOr` returns the value, or `def` when the key is missing. A key stored with value `""` is **present**.
- `CopyKeys` copies each key that exists in `src` into `dst` and returns how many it copied.
- `ReadOnly` returns a Getter that reads through to `g` (later changes to `g` are visible) but has **no other methods**, so asserting it to `Setter` or `Store` fails.
