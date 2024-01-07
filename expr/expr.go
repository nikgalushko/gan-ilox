package expr

import (
	"github.com/nikgalushko/gan-ilox/token"
)

type Expr interface {
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitBinaryExpr(expr Binary) any
	VisitGroupingExpr(expr Grouping) any
	VisitLiteralExpr(expr Literal) any
	VisitUnaryExpr(expr Unary) any
}

type StmtVisitor interface {
	VisitStmtExpression(expr StmtExpression) any
	VisitPrintStmt(s PrintStmt) any
}

type Binary struct {
	Left     Expr
	Operator token.Token
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

type Literal struct {
	Value token.Literal
}

func (e Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteralExpr(e)
}

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (e Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnaryExpr(e)
}

type Stmt interface {
	Accept(StmtVisitor) any
}

type StmtExpression struct {
	Expression Expr
}

func (e StmtExpression) Accept(v StmtVisitor) any {
	return v.VisitStmtExpression(e)
}

type PrintStmt struct {
	Expression Expr
}

func (e PrintStmt) Accept(v StmtVisitor) any {
	return v.VisitPrintStmt(e)
}
