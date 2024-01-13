package token

import (
	"fmt"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token/kind"
)

type Token struct {
	Type    kind.TokenType
	Lexeme  string
	Line    int
	Literal internal.Literal
}

func New(_type kind.TokenType, lexeme string, line int, l internal.Literal) Token {
	return Token{
		Type:    _type,
		Lexeme:  lexeme,
		Line:    line,
		Literal: l,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("Type: %+v; Literal: %+v", t.Type, t.Literal)
}
