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

	if !(left.IsFloat() || left.IsInt()) || !(right.IsFloat() || right.IsInt()) {
		i.err = errors.New("Type missmatch")
		return token.LiteralNil
	}

	switch expression.Operator.Kind {
	case token.Minus:
		if left.IsInt() && right.IsInt() {
			return token.NewLiteralInt(left.AsInt() - right.AsInt())
		} else {
			return token.NewLiteralFloat(left.AsFloat() - right.AsFloat())
		}
	case token.Plus:
		if left.IsInt() && right.IsInt() {
			return token.NewLiteralInt(left.AsInt() + right.AsInt())
		} else {
			return token.NewLiteralFloat(left.AsFloat() + right.AsFloat())
		}
	case token.Slash:
		if left.IsInt() && right.IsInt() {
			return token.NewLiteralInt(left.AsInt() / right.AsInt())
		} else {
			return token.NewLiteralFloat(left.AsFloat() / right.AsFloat())
		}
	case token.Star:
		if left.IsInt() && right.IsInt() {
			return token.NewLiteralInt(left.AsInt() * right.AsInt())
		} else {
			return token.NewLiteralFloat(left.AsFloat() * right.AsFloat())
		}
	}

	panic("unreachable code")
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

	switch expression.Operator.Kind {
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
