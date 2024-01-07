package interpreter

import (
	"errors"
	"fmt"

	"github.com/nikgalushko/gan-ilox/env"
	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Interpreter struct {
	env   *env.Environment
	stmts []expr.Stmt
	err   error
}

func New(env *env.Environment, stmts []expr.Stmt) *Interpreter {
	return &Interpreter{env: env, stmts: stmts}
}

func (i *Interpreter) Interpret() ([]any, error) {
	var ret []any
	for _, s := range i.stmts {
		v, err := i.exec(s)
		if err != nil {
			return nil, err
		}
		if v != nil && v.(token.Literal) != token.LiteralNil {
			ret = append(ret, v)
		}
	}

	return ret, nil
}

func (i *Interpreter) eval(e expr.Expr) (any, error) {
	ret := e.Accept(i)
	return ret, i.err
}

func (i *Interpreter) exec(s expr.Stmt) (any, error) {
	ret := s.Accept(i)
	return ret, i.err
}

func (i *Interpreter) VisitVarStmt(s expr.VarStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	value := token.LiteralNil
	name := s.Name.Lexeme

	if s.Expression != nil {
		v, err := i.eval(s.Expression)
		if err == nil {
			value = v.(token.Literal)
			i.err = err
		}
	}

	i.env.Set(name, value)

	return token.LiteralNil
}

func (i *Interpreter) VisitPrintStmt(s expr.PrintStmt) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	fmt.Println(val.(token.Literal).String())

	return token.LiteralNil
}

func (i *Interpreter) VisitStmtExpression(s expr.StmtExpression) any {
	if i.err != nil {
		return token.LiteralNil
	}

	ret, err := i.eval(s.Expression)
	if err != nil {
		i.err = err
		ret = token.LiteralNil
	}

	return ret
}

func (i *Interpreter) VisitAssignmentExpr(e expr.Assignment) any {
	if i.err != nil {
		return token.LiteralNil
	}

	if !i.env.Has(e.Name.Lexeme) {
		i.err = errors.New("undefined variable")
		return token.LiteralNil
	}

	val, err := i.eval(e.Expression)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	i.env.Set(e.Name.Lexeme, val.(token.Literal))

	return val
}

func (i *Interpreter) VisitVariableExpr(e expr.Variable) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val, err := i.env.Get(e.Name.Lexeme)
	if err != nil {
		i.err = err
		return token.LiteralNil
	}

	return val
}

func (i *Interpreter) VisitBinaryExpr(expression expr.Binary) any {
	if i.err != nil {
		return token.LiteralNil
	}

	left := expression.Left.Accept(i).(token.Literal)
	right := expression.Right.Accept(i).(token.Literal)

	var ret token.Literal
	switch expression.Operator.Type {
	case token.Minus:
		ret, i.err = sub(left, right)
	case token.Plus:
		ret, i.err = add(left, right)
	case token.Slash:
		ret, i.err = div(left, right)
	case token.Star:
		ret, i.err = mul(left, right)
	case token.Less:
		ret, i.err = less(left, right)
	case token.LessEqual:
		ret, i.err = lessOrEqual(left, right)
	case token.Greater:
		ret, i.err = graeater(left, right)
	case token.GreaterEqual:
		ret, i.err = graeaterOrEqual(left, right)
	case token.EqualEqual:
		ret, i.err = equal(left, right)
	}

	return ret
}

func (i *Interpreter) VisitGroupingExpr(expression expr.Grouping) any {
	if i.err != nil {
		return token.LiteralNil
	}
	return expression.Expression.Accept(i)
}

func (i *Interpreter) VisitLiteralExpr(expression expr.Literal) any {
	if i.err != nil {
		return token.LiteralNil
	}
	return expression.Value
}

func (i *Interpreter) VisitUnaryExpr(expression expr.Unary) any {
	if i.err != nil {
		return token.LiteralNil
	}

	val := expression.Right.Accept(i).(token.Literal)

	switch expression.Operator.Type {
	case token.Bang:
		return token.NewLiteralBool(!val.AsBool())
	case token.Minus:
		if val.IsInt() {
			return token.NewLiteralInt(-val.AsInt())
		} else if val.IsFloat() {
			return token.NewLiteralFloat(-val.AsFloat())
		}

		i.err = errors.New("Illegal operation") // TODO: craete more freandly error message
		return token.LiteralNil
	case token.BitwiseNot:
		if val.IsInt() {
			return token.NewLiteralInt(^val.AsInt())
		}
		i.err = errors.New("bitwise operator can be used only with integer number")
		return token.LiteralNil
	}

	panic("unreachable code")
}
