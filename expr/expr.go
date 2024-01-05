package expr

import (
	"github.com/nikgalushko/gan-ilox/token"
)

type Expr[R any] interface {
	Accept(visitor Visitor[R]) R
}

type Visitor[R any] interface {
	VisitBinaryExpr(expr Binary[R]) R
	VisitGroupingExpr(expr Grouping[R]) R
	VisitLiteralExpr(expr Literal[R]) R
	VisitUnaryExpr(expr Unary[R]) R
}

type Binary[R any] struct {
	Left     Expr[R]
	Operator token.Token
	Right    Expr[R]
}

func (e Binary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitBinaryExpr(e)
}

type Grouping[R any] struct {
	Expression Expr[R]
}

func (e Grouping[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitGroupingExpr(e)
}

type Literal[R any] struct {
	Value any
}

func (e Literal[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitLiteralExpr(e)
}

type Unary[R any] struct {
	Operator token.Token
	Right    Expr[R]
}

func (e Unary[R]) Accept(visitor Visitor[R]) R {
	return visitor.VisitUnaryExpr(e)
}
