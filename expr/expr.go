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
	VisitVariableExpr(expr Variable) any
	VisitAssignmentExpr(expr Assignment) any
}

type StmtVisitor interface {
	VisitStmtExpression(expr StmtExpression) any
	VisitPrintStmt(s PrintStmt) any
	VisitVarStmt(s VarStmt) any
	VisitBlockStmt(s BlockStmt) any
	VisitIfStmt(s IfStmt) any
	VisitElseStmt(s ElseStmt) any
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

type Variable struct {
	Name token.Token
}

func (e Variable) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariableExpr(e)
}

type Assignment struct {
	Name       token.Token
	Expression Expr
}

func (e Assignment) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssignmentExpr(e)
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

type VarStmt struct {
	Name       token.Token
	Expression Expr
}

func (e VarStmt) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(e)
}

type BlockStmt struct {
	Stmts []Stmt
}

func (e BlockStmt) Accept(v StmtVisitor) any {
	return v.VisitBlockStmt(e)
}

type IfStmt struct {
	Condition Expr
	If        Stmt
	Else      Stmt
}

func (e IfStmt) Accept(v StmtVisitor) any {
	return v.VisitIfStmt(e)
}

type ElseStmt struct {
	If    Stmt
	Block Stmt
}

func (e ElseStmt) Accept(v StmtVisitor) any {
	return v.VisitElseStmt(e)
}
