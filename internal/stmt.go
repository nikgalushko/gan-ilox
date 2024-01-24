package internal

type StmtVisitor interface {
	VisitStmtExpression(expr StmtExpression) any
	VisitPrintStmt(s PrintStmt) any
	VisitVarStmt(s VarStmt) any
	VisitBlockStmt(s BlockStmt) any
	VisitIfStmt(s IfStmt) any
	VisitElseStmt(s ElseStmt) any
	VisitForSmt(s ForStmt) any
	VisitFuncStmt(s FuncStmt) any
	VisitReturnStmt(s RreturnStmt) any
	VisitClassStmt(s ClassStmt) any
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
	Name       string
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

type ForStmt struct {
	Initializer Stmt
	Condition   Expr
	Step        Expr
	Body        Stmt
}

func (e ForStmt) Accept(v StmtVisitor) any {
	return v.VisitForSmt(e)
}

type FuncStmt struct {
	Name       string
	Parameters []string
	Body       Stmt
}

func (e FuncStmt) Accept(v StmtVisitor) any {
	return v.VisitFuncStmt(e)
}

type RreturnStmt struct {
	Expression Expr
}

func (e RreturnStmt) Accept(v StmtVisitor) any {
	return v.VisitReturnStmt(e)
}

type ClassStmt struct {
	Name    string
	Methods []FuncStmt
}

func (e ClassStmt) Accept(v StmtVisitor) any {
	return v.VisitClassStmt(e)
}
