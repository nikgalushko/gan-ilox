## Scanner
- [ ] fuzz tests
- [X] skip whitespaces
- [X] /* */ comment style

## Token
- [ ] replace `any` as type of literal to something like
```go
type ObjectKind int8
const (
    _ ObjectKind = iota
    FloatObject
    IntObject
    StringObject
)
type Object struct {
    i int64
    f float64
    s string
    kind ObjectKind
}
```