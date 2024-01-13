package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token"
	"github.com/nikgalushko/gan-ilox/token/kind"
)

type SyntaxError struct {
	lineNumber     int
	message, cause string
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("%d\t|\tError%s: %s", e.lineNumber, e.cause, e.message)
}

type Scanner struct {
	source               []rune
	start, current, line int
	tokens               []token.Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{
		source: []rune(source),
		line:   1,
	}
}

func (s *Scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, token.New(kind.EOF, "", s.line, internal.LiteralNil))

	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	r := s.advance()
	switch r {
	case '(':
		s.appendSingleToken(kind.LeftParen)
	case ')':
		s.appendSingleToken(kind.RightParen)
	case '{':
		s.appendSingleToken(kind.LeftBrace)
	case '}':
		s.appendSingleToken(kind.RightBrace)
	case ',':
		s.appendSingleToken(kind.Comma)
	case '.':
		s.appendSingleToken(kind.Dot)
	case ';':
		s.appendSingleToken(kind.Semicolon)
	case '+':
		s.appendSingleToken(kind.Plus)
	case '-':
		s.appendSingleToken(kind.Minus)
	case '*':
		s.appendSingleToken(kind.Star)
	case '&':
		s.appendSingleToken(kind.BitwiseAnd)
	case '|':
		s.appendSingleToken(kind.BitwiseOr)
	case '^':
		s.appendSingleToken(kind.BitwiseXor)
	case '~':
		s.appendSingleToken(kind.BitwiseNot)
	case '!':
		k := kind.Bang
		if s.match('=') {
			k = kind.BangEqual
		}
		s.appendSingleToken(k)
	case '=':
		k := kind.Equal
		if s.match('=') {
			k = kind.EqualEqual
		}
		s.appendSingleToken(k)
	case '<':
		k := kind.Less
		if s.match('=') {
			k = kind.LessEqual
		}
		s.appendSingleToken(k)
	case '>':
		k := kind.Greater
		if s.match('=') {
			k = kind.GreaterEqual
		}
		s.appendSingleToken(k)
	case '/':
		if s.match('/') {
			// skeep comment line
			for s.peek() != '\n' && !s.isAtEnd() {
				_ = s.advance()
			}
		} else if s.match('*') {
			var prevRune rune = 1
			for !(prevRune == '*' && s.peek() == '/') && !s.isAtEnd() {
				prevRune = s.advance()
				if prevRune == '\n' {
					s.line++
				}
			}
			_ = s.advance() // read last /
		} else {
			s.appendSingleToken(kind.Slash)
		}
	case ' ', '\r', '\t':
	case '\n':
		s.line++
	case '"':
		err := s.string()
		if err != nil {
			return err
		}
	default:
		if unicode.IsDigit(r) {
			if err := s.number(); err != nil {
				return err
			}
		} else if unicode.IsLetter(r) {
			if err := s.identifier(); err != nil {
				return err
			}
		} else {
			return SyntaxError{lineNumber: s.line, message: "Unexpected character"}
		}
	}

	return nil
}

func (s *Scanner) identifier() error {
	for unicode.IsDigit(s.peek()) || unicode.IsLetter(s.peek()) {
		_ = s.advance()
	}

	text := string(s.source[s.start:s.current])
	_type, ok := keywords[text]
	if !ok {
		_type = kind.Identifier
	}

	s.appendSingleToken(_type)
	return nil
}

func (s *Scanner) number() error {
	for unicode.IsDigit(s.peek()) {
		_ = s.advance()
	}

	isFloat := false
	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		isFloat = true
		_ = s.advance()

		for unicode.IsDigit(s.peek()) {
			_ = s.advance()
		}
	}

	text := string(s.source[s.start:s.current])
	var l internal.Literal

	if isFloat {
		n, err := strconv.ParseFloat(text, 64)
		if err != nil {
			return err
		}
		l = internal.NewLiteralFloat(n)
	} else {
		n, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return err
		}
		l = internal.NewLiteralInt(n)
	}

	s.tokens = append(s.tokens, token.New(kind.Number, text, s.line, l))
	return nil
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		_ = s.advance()
	}

	if s.isAtEnd() {
		return SyntaxError{lineNumber: s.line, message: "Untermintaed string"}
	}

	_ = s.advance()

	text := string(s.source[s.start+1 : s.current-1])
	s.tokens = append(s.tokens, token.New(kind.String, text, s.line, internal.NewLiteralString(text)))
	return nil
}

func (s *Scanner) appendSingleToken(_type kind.TokenType) {
	s.tokens = append(s.tokens, token.New(_type, string(s.source[s.start:s.current]), s.line, internal.LiteralNil))
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++

	return r
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return -1
	}

	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.isAtEnd() || s.current+1 >= len(s.source) {
		return -1
	}

	return s.source[s.current+1]
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

var keywords = map[string]kind.TokenType{
	"and":    kind.And,
	"or":     kind.Or,
	"class":  kind.Class,
	"if":     kind.If,
	"else":   kind.Else,
	"false":  kind.False,
	"true":   kind.True,
	"for":    kind.For,
	"while":  kind.While,
	"fun":    kind.Fun,
	"super":  kind.Super,
	"this":   kind.This,
	"print":  kind.Print,
	"return": kind.Return,
	"var":    kind.Var,
	"nil":    kind.Nil,
}
