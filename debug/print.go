package debug

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

type AstPrinter struct {
	E expr.Expr
	S []expr.Stmt
}

func (p AstPrinter) String() string {
	if len(p.S) == 0 {
		return p.E.Accept(p).(string)
	}

	var ret []string
	for _, s := range p.S {
		ret = append(ret, s.Accept(p).(string))
	}

	return strings.Join(ret, "\n")
}

func (p AstPrinter) VisitVarStmt(s expr.VarStmt) any {
	return p.parenthesize(s.Name.Lexeme, s.Expression)
}

func (p AstPrinter) VisitPrintStmt(s expr.PrintStmt) any {
	return p.parenthesize("print", s.Expression)
}

// TODO: rename to ExpressionStmt
func (p AstPrinter) VisitStmtExpression(s expr.StmtExpression) any {
	return p.parenthesize("stmt", s.Expression)
}

func (p AstPrinter) VisitAssignmentExpr(e expr.Assignment) any {
	return p.parenthesize(e.Name.Lexeme, e.Expression)
}

func (p AstPrinter) VisitVariableExpr(e expr.Variable) any {
	return e.Name.Lexeme
}

func (p AstPrinter) VisitBinaryExpr(expression expr.Binary) any {
	return p.parenthesize(expression.Operator.Lexeme, expression.Left, expression.Right)
}

func (p AstPrinter) VisitGroupingExpr(expression expr.Grouping) any {
	return p.parenthesize("group", expression.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expression expr.Literal) any {
	if expression.Value == token.LiteralNil {
		return "nil"
	}

	if expression.Value.IsInt() {
		return fmt.Sprintf("%d", expression.Value.AsInt())
	} else if expression.Value.IsFloat() {
		return fmt.Sprintf("%f", expression.Value.AsFloat())
	} else if expression.Value.IsBool() {
		return fmt.Sprintf("%t", expression.Value.AsBool())
	}

	return expression.Value.AsString()
}

func (p AstPrinter) VisitUnaryExpr(expression expr.Unary) any {
	return p.parenthesize(expression.Operator.Lexeme, expression.Right)
}

func (p AstPrinter) parenthesize(name string, expressions ...expr.Expr) any {
	out := bytes.NewBuffer(nil)
	fmt.Fprintf(out, "(%s", name)

	for _, e := range expressions {
		fmt.Fprintf(out, " ")
		fmt.Fprint(out, e.Accept(p))
	}
	fmt.Fprint(out, ")")

	return out.String()
}
