package expr

import (
	"github.com/nikgalushko/gan-ilox/token"
)

type Expr interface {
	Accept(visitor Visitor) any
}

type Visitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (e Binary) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e Grouping) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(e)
}

type Literal struct {
	Value any
}

func (e Literal) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e Unary) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(e)
}
