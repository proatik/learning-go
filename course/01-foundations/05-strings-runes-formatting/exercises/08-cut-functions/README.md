# Cut, CutPrefix and CutLast

Practise the three `Cut` functions. Implement in `cut.go`:

```go
func ParseKeyValue(s string) (key, value string, err error)
func StripScheme(url string) (rest string, secure bool)
func SplitExt(filename string) (name, ext string)
```

### `ParseKeyValue` — `strings.Cut`
Split at the **first** `=`, then trim surrounding whitespace from key and value.

| Input | key | value | err |
|---|---|---|---|
| `"name=Go"` | `"name"` | `"Go"` | `nil` |
| `"  lang = Go 1.27 "` | `"lang"` | `"Go 1.27"` | `nil` |
| `"url=a=b"` | `"url"` | `"a=b"` | `nil` |
| `"empty="` | `"empty"` | `""` | `nil` |
| `"no separator"` | `""` | `""` | `missing '=' in "no separator"` |
| `" = value"` | `""` | `""` | `empty key in " = value"` |

### `StripScheme` — `strings.CutPrefix`

| Input | rest | secure |
|---|---|---|
| `"https://go.dev"` | `"go.dev"` | `true` |
| `"http://example.com"` | `"example.com"` | `false` |
| `"go.dev"` | `"go.dev"` | `false` |

### `SplitExt` — `strings.CutLast` (new in Go 1.27)
Split at the **last** `.`.

| Input | name | ext |
|---|---|---|
| `"archive.tar.gz"` | `"archive.tar"` | `"gz"` |
| `"main.go"` | `"main"` | `"go"` |
| `"README"` | `"README"` | `""` |

## Rules
- On error, `ParseKeyValue` returns empty `key` and `value` and an error made with `fmt.Errorf` — the messages must match the table exactly (use `%q`).
- `StripScheme` removes `"https://"` (secure) or `"http://"` (not secure); anything else is returned unchanged.
