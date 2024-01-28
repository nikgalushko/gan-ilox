package internal

import (
	"errors"
	"strconv"
)

type Interpreter interface {
	Exec(Stmt[Literal]) (Literal, error)
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
	literalClass
	literalClassInstance
)

type Literal struct {
	i              int64
	f              float64
	s              string
	b              bool
	function       Function
	class          Class
	instance       ClassInstance
	_type          literalType
	isReturnResult bool
}

type ClassInstance struct {
	Class  *Class
	Fields map[string]Literal
}

func (c ClassInstance) Set(name string, value Literal) {
	c.Fields[name] = value
}

func (c ClassInstance) Get(name string) (Literal, error) {
	v, ok := c.Fields[name]
	if !ok {
		v, ok = c.Class.Methods[name]
		if !ok {
			return LiteralNil, errors.New("undefined field: " + name)
		}
	}

	return v, nil
}

type Function struct {
	ArgumentsName []string
	body          Stmt[Literal]
	f             func(args ...Literal) (Literal, error)
}

func (f Function) Call(params []Literal, i Interpreter) (Literal, error) {
	if f.f != nil {
		return f.f(params...)
	}

	return i.Exec(f.body)
}

type Class struct {
	Name        string
	Initializer *Literal
	Methods     map[string]Literal
}

func (c Class) Call(params []Literal, i Interpreter) (Literal, error) {
	if c.Initializer != nil {
		_, err := i.Exec(c.Initializer.function.body)
		if err != nil {
			return LiteralNil, err
		}
	}

	return Literal{_type: literalClassInstance, instance: ClassInstance{Class: &c, Fields: map[string]Literal{}}}, nil
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

func NewLiteralUserFunction(args []string, body Stmt[Literal]) Literal {
	return Literal{_type: literalFunction, function: Function{ArgumentsName: args, body: body}}
}

func NewLiteralNativeFunction(args []string, f func(args ...Literal) (Literal, error)) Literal {
	return Literal{_type: literalFunction, function: Function{ArgumentsName: args, f: f}}
}

func NewLiteralClass(name string, methods map[string]Literal) Literal {
	return Literal{_type: literalClass, class: Class{Name: name, Methods: methods}}
}

func (l Literal) IsClass() bool {
	return l._type == literalClass
}

func (l Literal) IsClassInstance() bool {
	return l._type == literalClassInstance
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

func (l Literal) AsClass() Class {
	return l.class
}

func (l Literal) AsClassInstance() ClassInstance {
	return l.instance
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
