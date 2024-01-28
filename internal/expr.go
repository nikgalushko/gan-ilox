package internal

import (
	"github.com/nikgalushko/gan-ilox/token/kind"
)

type Expr[In, Out any] interface {
	Accept(visitor ExprVisitor[In, Out]) Out
}

type ExprVisitor[In, Out any] interface {
	VisitBinaryExpr(expr Binary[In, Out]) Out
	VisitGroupingExpr(expr Grouping[In, Out]) Out
	VisitLiteralExpr(expr LiteralExpr[In, Out]) Out
	VisitUnaryExpr(expr Unary[In, Out]) Out
	VisitVariableExpr(expr Variable[In, Out]) Out
	VisitAssignmentExpr(expr Assignment[In, Out]) Out
	VisitLogicalExpr(e Logical[In, Out]) Out
	VisitCallExpr(e Call[In, Out]) Out
	VisitGetExpr(e GetExpr[In, Out]) Out
	VisitSetExpr(e SetExpr[In, Out]) Out
}

type Call[In, Out any] struct {
	Arguments []Expr[In, Out]
	Callee    Expr[In, Out]
}

func (e Call[In, Out]) Accept(v ExprVisitor[In, Out]) Out {
	return v.VisitCallExpr(e)
}

type Binary[In, Out any] struct {
	Left     Expr[In, Out]
	Operator kind.TokenType
	Right    Expr[In, Out]
}

func (e Binary[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitBinaryExpr(e)
}

type Grouping[In, Out any] struct {
	Expression Expr[In, Out]
}

func (e Grouping[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr[In, Out any] struct {
	Value Literal
}

func (e LiteralExpr[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitLiteralExpr(e)
}

type Unary[In, Out any] struct {
	Operator kind.TokenType
	Right    Expr[In, Out]
}

func (e Unary[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitUnaryExpr(e)
}

type Variable[In, Out any] struct {
	Name string
}

func (e Variable[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitVariableExpr(e)
}

type Assignment[In, Out any] struct {
	Name       string
	Expression Expr[In, Out]
}

func (e Assignment[In, Out]) Accept(visitor ExprVisitor[In, Out]) Out {
	return visitor.VisitAssignmentExpr(e)
}

type Logical[In, Out any] struct {
	Left     Expr[In, Out]
	Operator kind.TokenType
	Right    Expr[In, Out]
}

func (e Logical[In, Out]) Accept(v ExprVisitor[In, Out]) Out {
	return v.VisitLogicalExpr(e)
}

type GetExpr[In, Out any] struct {
	Name       string
	Expression Expr[In, Out]
}

func (e GetExpr[In, Out]) Accept(v ExprVisitor[In, Out]) Out {
	return v.VisitGetExpr(e)
}

type SetExpr[In, Out any] struct {
	Name   string
	Object Expr[In, Out]
	Value  Expr[In, Out]
}

func (e SetExpr[In, Out]) Accept(v ExprVisitor[In, Out]) Out {
	return v.VisitSetExpr(e)
}
