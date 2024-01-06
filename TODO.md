## Scanner
- [ ] fuzz tests
- [X] skip whitespaces
- [X] /* */ comment style

## Token
- [X] replace `any` as type of literal to something like
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

## Interpreter
- [ ] support runtime errors
- [ ] use `token.Literal` instead of `any`

## Expr
- [ ] rewrite codegen