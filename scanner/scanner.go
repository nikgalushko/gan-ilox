package scanner

import (
	"fmt"
	"strconv"
	"unicode"

	"github.com/nikgalushko/gan-ilox/token"
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

	s.tokens = append(s.tokens, token.New(token.EOF, "", s.line, nil))

	return s.tokens, nil
}

func (s *Scanner) scanToken() error {
	r := s.advance()
	switch r {
	case '(':
		s.appendSingleToken(token.LeftParen)
	case ')':
		s.appendSingleToken(token.RightParen)
	case '{':
		s.appendSingleToken(token.LeftBrace)
	case '}':
		s.appendSingleToken(token.RightBrace)
	case ',':
		s.appendSingleToken(token.Comma)
	case '.':
		s.appendSingleToken(token.Dot)
	case ';':
		s.appendSingleToken(token.Semicolon)
	case '+':
		s.appendSingleToken(token.Plus)
	case '-':
		s.appendSingleToken(token.Minus)
	case '*':
		s.appendSingleToken(token.Star)
	case '!':
		k := token.Bang
		if s.match('=') {
			k = token.BangEqual
		}
		s.appendSingleToken(k)
	case '=':
		k := token.Equal
		if s.match('=') {
			k = token.EqualEqual
		}
		s.appendSingleToken(k)
	case '<':
		k := token.Less
		if s.match('=') {
			k = token.LessEqual
		}
		s.appendSingleToken(k)
	case '>':
		k := token.Greater
		if s.match('=') {
			k = token.GreaterEqual
		}
		s.appendSingleToken(k)
	case '/':
		if s.match('/') {
			// skeep comment line
			for s.peek() != '\n' && !s.isAtEnd() {
				_ = s.advance()
			}
		} else {
			s.appendSingleToken(token.Slash)
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
	kind, ok := keywords[text]
	if !ok {
		kind = token.Identifier
	}

	s.appendSingleToken(kind)
	return nil
}

func (s *Scanner) number() error {
	for unicode.IsDigit(s.peek()) {
		_ = s.advance()
	}

	if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
		_ = s.advance()

		for unicode.IsDigit(s.peek()) {
			_ = s.advance()
		}
	}

	text := string(s.source[s.start:s.current])
	n, err := strconv.ParseFloat(text, 64)
	if err != nil {
		return err
	}

	s.tokens = append(s.tokens, token.New(token.Number, text, s.line, n))
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
	s.tokens = append(s.tokens, token.New(token.String, text, s.line, text))
	return nil
}

func (s *Scanner) appendSingleToken(kind token.TokenKind) {
	s.tokens = append(s.tokens, token.New(kind, string(s.source[s.start:s.current]), s.line, nil))
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

var keywords = map[string]token.TokenKind{
	"and":    token.And,
	"or":     token.Or,
	"class":  token.Class,
	"if":     token.If,
	"else":   token.Else,
	"false":  token.False,
	"true":   token.True,
	"for":    token.For,
	"while":  token.While,
	"fun":    token.Fun,
	"super":  token.Super,
	"this":   token.This,
	"print":  token.Print,
	"return": token.Return,
	"var":    token.Var,
	"nil":    token.Nil,
}
