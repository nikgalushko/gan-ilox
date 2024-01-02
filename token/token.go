package token

type TokenKind int

const (
	_ TokenKind = iota
	// single-character tokens
	LeftParen
	RightParen
	LeftBrace
	RightBrace
	Comma
	Dot
	Minus
	Plus
	Semicolon
	Slash
	Star

	// One or two character tokens
	Bang
	BangEqual
	Equal
	EqualEqual
	Greater
	GreaterEqual
	Less
	LessEqual

	// Literals
	Identifier
	String
	Number

	// Keywords
	And
	Class
	Else
	False
	Fun
	For
	If
	Nil
	Or
	Print
	Return
	Super
	This
	True
	Var
	While

	EOF
)

type Token struct {
	Kind    TokenKind
	Lexeme  string
	Line    int
	Literal any
}

func New(kind TokenKind, lexeme string, line int, literal any) Token {
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
