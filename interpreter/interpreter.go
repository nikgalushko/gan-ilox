package interpreter

import (
	"math"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type Interpreter struct {
	E expr.Expr
}

func New(E expr.Expr) Interpreter {
	return Interpreter{E: E}
}

func (i Interpreter) Eval() any {
	return i.E.Accept(i)
}

func (i Interpreter) VisitBinaryExpr(expression expr.Binary) any {
	left := expression.Left.Accept(i)
	right := expression.Right.Accept(i)

	switch expression.Operator.Kind {
	case token.Minus:
		return toNumber(left) - toNumber(right)
	case token.Plus:
		return toNumber(left) + toNumber(right)
	case token.Slash:
		return toNumber(left) / toNumber(right)
	case token.Star:
		return toNumber(left) * toNumber(right)
	}

	panic("unreachable code")
}

func (i Interpreter) VisitGroupingExpr(expression expr.Grouping) any {
	return expression.Expression.Accept(i)
}

func (i Interpreter) VisitLiteralExpr(expression expr.Literal) any {
	return expression.Value
}

func (i Interpreter) VisitUnaryExpr(expression expr.Unary) any {
	// TODO: check the type of val
	val := expression.Right.Accept(i)

	switch expression.Operator.Kind {
	case token.Bang:
		return !toBool(val)
	case token.Minus:
		return -toNumber(val)
	case token.BitwiseNot:
		return ^int64(toNumber(val))
	}

	panic("unreachable code")
}

func toNumber(v interface{}) float64 {
	switch n := v.(type) {
	case float64:
		return float64(n)
	case float32:
		return float64(n)
	case uint:
		return float64(n)
	case uint16:
		return float64(n)
	case uint32:
		return float64(n)
	case uint64:
		return float64(n)
	case int:
		return float64(n)
	case int16:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	}
	return math.NaN()
}

func toBool(val any) bool {
	if val == nil {
		return false
	}

	b, ok := val.(bool)
	if ok {
		return b
	}

	return true
}
