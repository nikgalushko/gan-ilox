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
	BitwiseAnd
	BitwiseOr
	BitwiseXor
	BitwiseNot

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

func (k TokenKind) String() string {
	switch k {
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case Comma:
		return ","
	case Dot:
		return "."
	case Minus:
		return "-"
	case Plus:
		return "+"
	case Semicolon:
		return ";"
	case Slash:
		return "/"
	case Star:
		return "*"

	// One or two character tokens
	case Bang:
		return "!"
	case BangEqual:
		return "!="
	case Equal:
		return "="
	case EqualEqual:
		return "=="
	case Greater:
		return ">"
	case GreaterEqual:
		return ">="
	case Less:
		return "<"
	case LessEqual:
		return "<="

	// Literals
	case Identifier:
		return "identifier"
	case String:
		return "string"
	case Number:
		return "number"

	// Keywords
	case And:
		return "and"
	case Class:
		return "class"
	case Else:
		return "else"
	case False:
		return "false"
	case Fun:
		return "fun"
	case For:
		return "for"
	case If:
		return "if"
	case Nil:
		return "nil"
	case Or:
		return "or"
	case Print:
		return "print"
	case Return:
		return "return"
	case Super:
		return "super"
	case This:
		return "this"
	case True:
		return "true"
	case Var:
		return "var"
	case While:
		return "while"

	case EOF:
		return "EOF"
	}

	return "<undefined>"
}

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
