package parser

import (
	"errors"
	"strings"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

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

func (p *Parser) Parse() ([]expr.Stmt, error) {
	var (
		pErr  PraseError
		stmts []expr.Stmt
	)

	for !p.isAtEnd() {
		s, err := p.declaration()
		if err != nil {
			pErr = append(pErr, err)
			p.synchronize()
		} else {
			stmts = append(stmts, s)
		}
	}

	if len(pErr) == 0 {
		return stmts, nil
	}

	return stmts, pErr
}

func (p *Parser) declaration() (expr.Stmt, error) {
	if p.match(token.Var) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() (expr.Stmt, error) {
	if !p.match(token.Identifier) {
		return nil, errors.New("expect variable name")
	}

	name := p.prev() // consume token in p.match
	var (
		initializer Expr
		err         error
	)

	if p.match(token.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(token.Semicolon) {
			return nil, errors.New("expect ; after variabl declaration")
		}
	}

	return expr.VarStmt{Name: name, Expression: initializer}, nil
}

func (p *Parser) statement() (expr.Stmt, error) {
	if p.match(token.Print) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (expr.Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(token.Semicolon) {
		return nil, errors.New("expected ; after expression")
	}

	return expr.PrintStmt{Expression: e}, nil
}

func (p *Parser) expressionStatement() (expr.Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(token.Semicolon) {
		return nil, errors.New("expected ; after expression")
	}
	return expr.StmtExpression{Expression: e}, nil
}

type Expr expr.Expr

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(token.Equal) {
		// e is variable
		variable, ok := e.(expr.Variable)
		if !ok {
			return nil, errors.New("invalid assignment target")
		}

		e, err := p.assignment()
		if err != nil {
			return nil, err
		}

		return expr.Assignment{Name: variable.Name, Expression: e}, nil
	}

	return e, nil
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
		e = expr.Binary{Left: e, Operator: operator, Right: right}
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
		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) term() (Expr, error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.Plus, token.Minus, token.BitwiseAnd, token.BitwiseOr) {
		operator := p.prev()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) factor() (Expr, error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.Slash, token.Star, token.BitwiseXor) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		e = expr.Binary{Left: e, Operator: operator, Right: right}
	}

	return e, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(token.Bang, token.Minus, token.BitwiseNot) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return expr.Unary{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(token.Number, token.String) {
		return expr.Literal{Value: p.prev().Literal}, nil
	}

	if p.match(token.True) {
		return expr.Literal{Value: token.NewLiteralBool(true)}, nil
	}
	if p.match(token.False) {
		return expr.Literal{Value: token.NewLiteralBool(false)}, nil
	}
	if p.match(token.Nil) {
		return expr.Literal{Value: token.LiteralNil}, nil
	}
	if p.match(token.Identifier) {
		return expr.Variable{Name: p.prev()}, nil
	}

	if p.match(token.LeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(token.RightParen) {
			return nil, errors.New("expect ')' after expression")
		}

		return expr.Grouping{Expression: e}, nil
	}

	return nil, errors.New("expect expression")
}

func (p *Parser) prev() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			_ = p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
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
	return p.peek().Type == token.EOF
}

func (p *Parser) synchronize() {
	t := p.advance()

	for !p.isAtEnd() {
		if p.prev().Type == token.Semicolon {
			return
		}

		switch t.Type {
		case token.Var, token.For, token.While, token.If, token.Else, token.Return, token.Print, token.Fun, token.Class:
			return
		}

		t = p.advance()
	}
}
