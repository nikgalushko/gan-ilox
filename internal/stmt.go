package internal

type StmtVisitor[In, Out any] interface {
	VisitStmtExpression(expr StmtExpression[In, Out]) Out
	VisitPrintStmt(s PrintStmt[In, Out]) Out
	VisitVarStmt(s VarStmt[In, Out]) Out
	VisitBlockStmt(s BlockStmt[In, Out]) Out
	VisitIfStmt(s IfStmt[In, Out]) Out
	VisitElseStmt(s ElseStmt[In, Out]) Out
	VisitForSmt(s ForStmt[In, Out]) Out
	VisitFuncStmt(s FuncStmt[In, Out]) Out
	VisitReturnStmt(s RreturnStmt[In, Out]) Out
	VisitClassStmt(s ClassStmt[In, Out]) Out
}

type Stmt[In, Out any] interface {
	Accept(StmtVisitor[In, Out]) Out
}

type StmtExpression[In, Out any] struct {
	Expression Expr[In, Out]
}

func (e StmtExpression[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitStmtExpression(e)
}

type PrintStmt[In, Out any] struct {
	Expression Expr[In, Out]
}

func (e PrintStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitPrintStmt(e)
}

type VarStmt[In, Out any] struct {
	Name       string
	Expression Expr[In, Out]
}

func (e VarStmt[In, Out]) Accept(visitor StmtVisitor[In, Out]) Out {
	return visitor.VisitVarStmt(e)
}

type BlockStmt[In, Out any] struct {
	Stmts []Stmt[In, Out]
}

func (e BlockStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitBlockStmt(e)
}

type IfStmt[In, Out any] struct {
	Condition Expr[In, Out]
	If        Stmt[In, Out]
	Else      Stmt[In, Out]
}

func (e IfStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitIfStmt(e)
}

type ElseStmt[In, Out any] struct {
	If    Stmt[In, Out]
	Block Stmt[In, Out]
}

func (e ElseStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitElseStmt(e)
}

type ForStmt[In, Out any] struct {
	Initializer Stmt[In, Out]
	Condition   Expr[In, Out]
	Step        Expr[In, Out]
	Body        Stmt[In, Out]
}

func (e ForStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitForSmt(e)
}

type FuncStmt[In, Out any] struct {
	Name       string
	Parameters []string
	Body       Stmt[In, Out]
}

func (e FuncStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitFuncStmt(e)
}

type RreturnStmt[In, Out any] struct {
	Expression Expr[In, Out]
}

func (e RreturnStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitReturnStmt(e)
}

type ClassStmt[In, Out any] struct {
	Name    string
	Methods []FuncStmt[In, Out]
}

func (e ClassStmt[In, Out]) Accept(v StmtVisitor[In, Out]) Out {
	return v.VisitClassStmt(e)
}
