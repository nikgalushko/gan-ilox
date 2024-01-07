package interpreter

import (
	"errors"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Interpreter struct {
	e   expr.Expr
	err error
}

func New(E expr.Expr) *Interpreter {
	return &Interpreter{e: E}
}

func (i *Interpreter) Eval() (any, error) {
	ret := i.e.Accept(i)
	return ret, i.err
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
