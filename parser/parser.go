package parser

import (
	"errors"
	"strings"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token"
	"github.com/nikgalushko/gan-ilox/token/kind"
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
	tokens         []token.Token
	current        int
	insideFunction bool
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) Parse() ([]internal.Stmt, error) {
	var (
		pErr  PraseError
		stmts []internal.Stmt
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

func (p *Parser) declaration() (internal.Stmt, error) {
	if p.match(kind.Var) {
		return p.varDeclaration()
	} else if p.match(kind.Fun) {
		return p.funDeclaration()
	}

	return p.statement()
}

func (p *Parser) funDeclaration() (internal.Stmt, error) {
	if !p.match(kind.Identifier) {
		return nil, errors.New("expect function name")
	}

	name := p.prev().Lexeme // consume token in p.match

	if !p.match(kind.LeftParen) {
		return nil, errors.New("expect '(' after function name")
	}

	ret := internal.FuncStmt{Name: name}
	defer func() { p.insideFunction = false }()

	if !p.match(kind.RightParen) {
		var args []string
		expectComma := false
		for !p.match(kind.RightParen) && !p.isAtEnd() {
			if expectComma && !p.match(kind.Comma) {
				return nil, errors.New("arguments must be splitted by comma")
			}

			if !p.match(kind.Identifier) {
				return nil, errors.New("expect argument name")
			}
			args = append(args, p.prev().Lexeme)
			expectComma = true
		}
		ret.Parameters = args
	}

	if !p.match(kind.LeftBrace) {
		return nil, errors.New("expect '{' as start of function body")
	}

	p.insideFunction = true
	body, err := p.blockStmt()
	if err != nil {
		return nil, err
	}

	ret.Body = body
	return ret, nil
}

func (p *Parser) varDeclaration() (internal.Stmt, error) {
	if !p.match(kind.Identifier) {
		return nil, errors.New("expect variable name")
	}

	name := p.prev() // consume token in p.match
	var (
		initializer Expr
		err         error
	)

	if p.match(kind.Equal) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(kind.Semicolon) {
			return nil, errors.New("expect ; after variabl declaration")
		}
	}

	return internal.VarStmt{Name: name.Lexeme, Expression: initializer}, nil
}

func (p *Parser) statement() (internal.Stmt, error) {
	if p.match(kind.Print) {
		return p.printStatement()
	} else if p.match(kind.LeftBrace) {
		return p.blockStmt()
	} else if p.match(kind.If) {
		return p.ifStmt()
	} else if p.match(kind.For) {
		return p.forStmt()
	} else if p.match(kind.Return) {
		return p.returnStmt()
	} else if p.match(kind.Class) {
		return p.classStmt()
	}
	return p.expressionStatement()
}

func (p *Parser) classStmt() (internal.Stmt, error) {
	if !p.match(kind.Identifier) {
		return nil, errors.New("expect class name")
	}

	name := p.prev().Lexeme // consume token in p.match

	if !p.match(kind.LeftBrace) {
		return nil, errors.New("expect '{' after class name")
	}

	var methods []internal.FuncStmt
	for !p.check(kind.RightBrace) && !p.isAtEnd() {
		m, err := p.funDeclaration()
		if err != nil {
			return nil, err
		}
		methods = append(methods, m.(internal.FuncStmt))
	}

	return internal.ClassStmt{Name: name, Methods: methods}, nil
}

func (p *Parser) returnStmt() (internal.Stmt, error) {
	if !p.insideFunction {
		return nil, errors.New("return statement is not inside a function")
	}
	ret := internal.RreturnStmt{}
	if !p.match(kind.Semicolon) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}
		ret.Expression = e
	}

	if !p.match(kind.Semicolon) {
		return ret, errors.New("expect ';' after return")
	}

	return ret, nil
}

func (p *Parser) forStmt() (internal.Stmt, error) {
	if !p.match(kind.LeftParen) {
		return nil, errors.New("expect '(' after for")
	}

	var (
		initializer internal.Stmt
		condition   internal.Expr
		err         error
	)
	if p.match(kind.Var) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else {
		v, err := p.expression()
		if err != nil {
			return nil, err
		}
		if p.match(kind.Semicolon) {
			initializer = internal.StmtExpression{Expression: v}
		} else {
			condition = v
		}
	}

	ret := internal.ForStmt{}
	if initializer == nil {
		ret.Condition = condition
	} else {
		ret.Initializer = initializer
		ret.Condition, err = p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(kind.Semicolon) {
			return nil, errors.New("expect ';' after for condition")
		}

		ret.Step, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(kind.RightParen) {
		return nil, errors.New("expect ')' as end for clauses")
	}

	if !p.match(kind.LeftBrace) {
		return nil, errors.New("expect '{' before for block")
	}

	ret.Body, err = p.blockStmt()

	return ret, err
}

func (p *Parser) ifStmt() (internal.Stmt, error) {
	if !p.match(kind.LeftParen) {
		return nil, errors.New("expect '(' after if")
	}

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(kind.RightParen) {
		return nil, errors.New("expect ')' after if condition")
	}

	if !p.match(kind.LeftBrace) {
		return nil, errors.New("expect '{' before if block")
	}

	ifBlock, err := p.blockStmt()
	if err != nil {
		return nil, err
	}

	ret := internal.IfStmt{Condition: condition, If: ifBlock}
	if p.match(kind.Else) {
		if p.match(kind.If) {
			ret.Else, err = p.ifStmt()
		} else if p.match(kind.LeftBrace) {
			ret.Else, err = p.blockStmt()
		} else {
			return nil, errors.New("unexpected symbol after else")
		}
	}

	return ret, err
}

// TODO: how to refactor this with Parse()
func (p *Parser) blockStmt() (internal.Stmt, error) {
	var (
		pErr  PraseError
		stmts []internal.Stmt
	)

	for !p.check(kind.RightBrace) && !p.isAtEnd() {
		s, err := p.declaration()
		if err != nil {
			pErr = append(pErr, err)
			p.synchronize()
		} else {
			stmts = append(stmts, s)
		}
	}

	if !p.match(kind.RightBrace) {
		pErr = append(pErr, errors.New("expect } after block"))
	}

	if len(pErr) == 0 {
		return internal.BlockStmt{Stmts: stmts}, nil
	}

	return nil, pErr
}

func (p *Parser) printStatement() (internal.Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(kind.Semicolon) {
		return nil, errors.New("expected ; after expression")
	}

	return internal.PrintStmt{Expression: e}, nil
}

func (p *Parser) expressionStatement() (internal.Stmt, error) {
	e, err := p.expression()
	if err != nil {
		return nil, err
	}

	if !p.match(kind.Semicolon) {
		return nil, errors.New("expected ; after expression")
	}
	return internal.StmtExpression{Expression: e}, nil
}

type Expr = internal.Expr

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	e, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(kind.Equal) {
		variable, ok := e.(internal.Variable)
		if !ok {
			return nil, errors.New("invalid assignment target")
		}

		e, err := p.assignment()
		if err != nil {
			return nil, err
		}

		return internal.Assignment{Name: variable.Name, Expression: e}, nil
	}

	return e, nil
}

func (p *Parser) or() (Expr, error) {
	e, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(kind.Or) {
		operator := p.prev()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		e = internal.Logical{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) and() (Expr, error) {
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(kind.And) {
		operator := p.prev()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		e = internal.Logical{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) equality() (Expr, error) {
	e, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(kind.EqualEqual, kind.BangEqual) {
		operator := p.prev()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		e = internal.Binary{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) comparison() (Expr, error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(kind.Less, kind.LessEqual, kind.Greater, kind.GreaterEqual) {
		operator := p.prev()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		e = internal.Binary{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) term() (Expr, error) {
	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(kind.Plus, kind.Minus, kind.BitwiseAnd, kind.BitwiseOr) {
		operator := p.prev()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		e = internal.Binary{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) factor() (Expr, error) {
	e, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(kind.Slash, kind.Star, kind.BitwiseXor) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		e = internal.Binary{Left: e, Operator: operator.Type, Right: right}
	}

	return e, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(kind.Bang, kind.Minus, kind.BitwiseNot) {
		operator := p.prev()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return internal.Unary{Operator: operator.Type, Right: right}, nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	e, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(kind.LeftParen) {
			e, err = p.finishCall(e)
			if err != nil {
				return nil, err
			}
		} else if p.match(kind.Dot) {
			if !p.match(kind.Identifier) {
				return nil, errors.New("expect property name after '.'")
			}
			e = internal.GetExpr{Name: p.prev().Lexeme, Expression: e}
		} else {
			break
		}
	}

	return e, nil
}

func (p *Parser) finishCall(callee Expr) (Expr, error) {
	var args []Expr
	if !p.check(kind.RightParen) {
		for {
			arg, err := p.expression()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)

			if !p.match(kind.Comma) {
				break
			}
		}
	}

	if !p.match(kind.RightParen) {
		return nil, errors.New("expect ')' as end of arguments")
	}

	return internal.Call{Callee: callee, Arguments: args}, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(kind.Number, kind.String) {
		return internal.LiteralExpr{Value: p.prev().Literal}, nil
	}

	if p.match(kind.True) {
		return internal.LiteralExpr{Value: internal.NewLiteralBool(true)}, nil
	}
	if p.match(kind.False) {
		return internal.LiteralExpr{Value: internal.NewLiteralBool(false)}, nil
	}
	if p.match(kind.Nil) {
		return internal.LiteralExpr{Value: internal.LiteralNil}, nil
	}
	if p.match(kind.Identifier) {
		return internal.Variable{Name: p.prev().Lexeme}, nil
	}

	if p.match(kind.LeftParen) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(kind.RightParen) {
			return nil, errors.New("expect ')' after expression")
		}

		return internal.Grouping{Expression: e}, nil
	}

	return nil, errors.New("expect expression")
}

func (p *Parser) prev() token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) match(tokens ...kind.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			_ = p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) check(t kind.TokenType) bool {
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
	return p.peek().Type == kind.EOF
}

func (p *Parser) synchronize() {
	t := p.advance()

	for !p.isAtEnd() {
		if p.prev().Type == kind.Semicolon {
			return
		}

		switch t.Type {
		case kind.Var, kind.For, kind.While, kind.If, kind.Else, kind.Return, kind.Print, kind.Fun, kind.Class:
			return
		}

		t = p.advance()
	}
}
