package internal

import (
	"github.com/nikgalushko/gan-ilox/token/kind"
)

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr LiteralExpr) any
	VisitUnaryExpr(expr Unary) any
	VisitVariableExpr(expr Variable) any
	VisitAssignmentExpr(expr Assignment) any
	VisitLogicalExpr(e Logical) any
	VisitCallExpr(e Call) any
	VisitGetExpr(e GetExpr) any
	VisitSetExpr(e SetExpr) any
}

type Call struct {
	Arguments []Expr
	Callee    Expr
}

func (e Call) Accept(v ExprVisitor) any {
	return v.VisitCallExpr(e)
}

type Binary struct {
	Left     Expr
	Operator kind.TokenType
	Right    Expr
}

func (e Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinaryExpr(e)
}

type Grouping struct {
	Expression Expr
}

func (e Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGroupingExpr(e)
}

type LiteralExpr struct {
	Value Literal
}

func (e LiteralExpr) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

type Unary struct {
	Operator kind.TokenType
	Right    Expr
}

func (e Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

type Variable struct {
	Name string
}

func (e Variable) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(e)
}

type Assignment struct {
	Name       string
	Expression Expr
}

func (e Assignment) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignmentExpr(e)
}

type Logical struct {
	Left     Expr
	Operator kind.TokenType
	Right    Expr
}

func (e Logical) Accept(v ExprVisitor) any {
	return v.VisitLogicalExpr(e)
}

type GetExpr struct {
	Name       string
	Expression Expr
}

func (e GetExpr) Accept(v ExprVisitor) any {
	return v.VisitGetExpr(e)
}

type SetExpr struct {
	Name   string
	Object Expr
	Value  Expr
}

func (e SetExpr) Accept(v ExprVisitor) any {
	return v.VisitSetExpr(e)
}
