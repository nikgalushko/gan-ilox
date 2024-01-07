package token

import "strconv"

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Line    int
	Literal Literal
}

func New(kind TokenKind, lexeme string, line int, literal Literal) Token {
	return Token{
		Kind:    kind,
		Lexeme:  lexeme,
		Line:    line,
		Literal: literal,
	}
}

func (t Token) String() string {
	panic("not implemented yet")
}

type literalType int8

const (
	_ literalType = iota
	literalNil
	literalFloat
	literalInt
	literalString
	literalBool
)

type Literal struct {
	i     int64
	f     float64
	s     string
	b     bool
	_type literalType
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
