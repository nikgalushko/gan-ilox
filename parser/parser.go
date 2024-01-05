package parser

import (
	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Expr = expr.Expr[any]

type Parser struct {
	tokens  []token.Token
	current int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	e := p.comparison()

	for p.match(token.EqualEqual, token.BangEqual) {
		operator := p.prev()
		right := p.comparison()
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e
}

func (p *Parser) comparison() Expr {
	e := p.term()

	for p.match(token.Less, token.LessEqual, token.Greater, token.GreaterEqual) {
		operator := p.prev()
		right := p.term()
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e
}

func (p *Parser) term() Expr {
	e := p.factor()

	for p.match(token.Plus, token.Minus) {
		operator := p.prev()
		right := p.factor()
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e
}

func (p *Parser) factor() Expr {
	e := p.unary()

	for p.match(token.Slash, token.Star) {
		operator := p.prev()
		right := p.unary()
		e = expr.Binary[any]{Left: e, Operator: operator, Right: right}
	}

	return e
}

func (p *Parser) unary() Expr {
	if p.match(token.Bang, token.Minus) {
		operator := p.prev()
		right := p.unary()
		return expr.Unary[any]{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(token.Number, token.String) {
		return expr.Literal[any]{Value: p.prev().Literal}
	}

	if p.match(token.True) {
		return expr.Literal[any]{Value: true}
	}
	if p.match(token.False) {
		return expr.Literal[any]{Value: false}
	}
	if p.match(token.Nil) {
		return expr.Literal[any]{Value: nil}
	}

	if p.match(token.LeftParen) {
		e := p.expression()
		if !p.match(token.RightParen) {
			panic("Expect ')' after expression")
		}

		return expr.Grouping[any]{Expression: e}
	}

	panic("unreachable code")
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
