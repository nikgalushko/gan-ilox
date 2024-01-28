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

### Parser
- [ ] unit tests

## Interpreter
- [ ] support runtime errors
- [X] use `internal.Literal` instead of `any`
- [ ] more frendly error message
- [ ] unit tests

## Expr
- [ ] rewrite codegen
- [ ] add to each expression type `Expr` suffix

### Language
- [ ] support `break`, `continue` in for-loop
- [ ] support concatenation between string and number
- [ ] support array and slice
- [X] infinite loop
- [ ] fix grammar
    - [ ] ifStatement contains additional `;`
    - [ ] forStatement contains additional `;`
- [ ] limit number of function argument
- [ ] check count of arguments and number of funtion parameters; in parsing time not in runtime
- [X] match arguments with parameters in function call
- [X] return statement
- [X] check closure
- [ ] anonymous functions
- [ ] semantic analysis
- [ ] support `this`
- [ ] pass arguments to constructor
- [ ] custom constructor
- [ ] set each instance uniq identifier
- [ ] support for-loop with empty init action `for (i < 10; i = i + 2)`
- [ ] support for-loop with empty condition action `for (var i = 0;; i = i + 2)`
- [ ] support for-loop with empty step action `for (var i = 0; i < 10)`
