package parser

import (
	"errors"
	"strings"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Expr = expr.Expr[any]

type PraseError []error

func (e PraseError) Error() string {
	var arr []string
	for _, err := range e {
		arr = append(arr, err.Error())
	}

	return strings.Join(arr, "\n")
}

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() (Expr, error) {
	var pErr PraseError
	for !p.isAtEnd() {
		e, err := p.expression()
		if err == nil || p.isAtEnd() {
			return e, err
		}

		pErr = append(pErr, err)
		p.synchronize()
	}

	return nil, pErr
}

func (p *Parser) synchronize() {
	t := p.advance()

	for !p.isAtEnd() {
		if p.prev().Kind == token.Semicolon {
			return
		}

		switch t.Kind {
		case token.Var, token.For, token.While, token.If, token.Else, token.Return, token.Print, token.Fun, token.Class:
			return
		}

		t = p.advance()
	}
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.EqualEqual, token.BangEqual) {
		operator := p.prev()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) comparison() (Expr, error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.Less, token.LessEqual, token.Greater, token.GreaterEqual) {
		operator := p.prev()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) term() (Expr, error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.Plus, token.Minus) {
		operator := p.prev()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) factor() (Expr, error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.Slash, token.Star) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(token.Bang, token.Minus) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expr.Unary[any]{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(token.Number, token.String) {
		return expr.Literal[any]{Value: p.prev().Literal}, nil
	}

	if p.match(token.True) {
		return expr.Literal[any]{Value: true}, nil
	}
	if p.match(token.False) {
		return expr.Literal[any]{Value: false}, nil
	}
	if p.match(token.Nil) {
		return expr.Literal[any]{Value: nil}, nil
	}

	if p.match(token.LeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(token.RightParen) {
			panic("Expect ')' after expression")
		}

		return expr.Grouping[any]{Expression: e}, nil
	}

	return nil, errors.New("Expect expression")
}

func (p *Parser) prev() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tokens ...token.TokenKind) bool {
	for _, t := range tokens {
		if p.check(t) {
			_ = p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.TokenKind) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Kind == t
}

func (p *Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p *Parser) advance() token.Token {
	if p.isAtEnd() {
		return p.prev()
	}

	ret := p.tokens[p.current]
	p.current++

	return ret
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}
