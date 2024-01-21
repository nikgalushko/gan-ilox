package internal

import (
	"strconv"
)

type Interpreter interface {
	Exec(Stmt) (any, error)
}

type literalType int8

const (
	_ literalType = iota
	literalNil
	literalFloat
	literalInt
	literalString
	literalBool
	literalFunction
)

type Literal struct {
	i              int64
	f              float64
	s              string
	b              bool
	function       Function
	_type          literalType
	isReturnResult bool
}

type Function struct {
	ArgumentsName []string
	body          Stmt
	f             func(args ...Literal) (Literal, error)
}

func (f Function) Call(params []Literal, i Interpreter) (any, error) {
	if f.f != nil {
		return f.f(params...)
	}

	return i.Exec(f.body)
}

var LiteralNil = Literal{_type: literalNil}

func NewLiteralInt(i int64) Literal {
	return Literal{i: i, _type: literalInt}
}

func NewLiteralFloat(f float64) Literal {
	return Literal{f: f, _type: literalFloat}
}

func NewLiteralBool(b bool) Literal {
	return Literal{b: b, _type: literalBool}
}

func NewLiteralString(s string) Literal {
	return Literal{s: s, _type: literalString}
}

func NewLiteralUserFunction(args []string, body Stmt) Literal {
	return Literal{_type: literalFunction, function: Function{ArgumentsName: args, body: body}}
}

func NewLiteralNativeFunction(args []string, f func(args ...Literal) (Literal, error)) Literal {
	return Literal{_type: literalFunction, function: Function{ArgumentsName: args, f: f}}
}

func (l Literal) IsFunction() bool {
	return l._type == literalFunction
}

func (l Literal) IsNumber() bool {
	return l.IsInt() || l.IsFloat()
}

func (l Literal) IsInt() bool {
	return l._type == literalInt
}

func (l Literal) IsFloat() bool {
	return l._type == literalFloat
}

func (l Literal) IsString() bool {
	return l._type == literalString
}

func (l Literal) IsBool() bool {
	return l._type == literalBool
}

func (l Literal) IsNil() bool {
	return l._type == literalNil
}

func (l Literal) IsReturnResult() bool {
	return l.isReturnResult
}

func (l Literal) AsInt() int64 {
	if l._type == literalFloat {
		return int64(l.f)
	}
	return l.i
}

func (l Literal) AsFloat() float64 {
	if l._type == literalInt {
		return float64(l.i)
	}
	return l.f
}

func (l Literal) AsString() string {
	return l.s
}

func (l Literal) AsBool() bool {
	if l.IsNil() {
		return false
	}

	if l.IsBool() {
		return l.b
	}

	return true
}

func (l Literal) AsFunction() Function {
	return l.function
}

func (l Literal) AsReturnResult() Literal {
	l.isReturnResult = true
	return l
}

func (l Literal) String() string {
	var ret string
	if l.IsInt() {
		ret = strconv.FormatInt(l.i, 10)
	} else if l.IsFloat() {
		ret = strconv.FormatFloat(l.f, 'e', 10, 64)
	} else if l.IsBool() {
		ret = strconv.FormatBool(l.b)
	} else if l.IsString() {
		ret = l.s
	} else {
		ret = "nil"
	}

	return ret
}
