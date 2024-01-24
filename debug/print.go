package debug

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/nikgalushko/gan-ilox/internal"
)

type AstPrinter struct {
	E internal.Expr
	S []internal.Stmt
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

func (p AstPrinter) VisitFuncStmt(s internal.FuncStmt) any {
	ret := []string{
		"(func " + s.Name + "(" + strings.Join(s.Parameters, ",") + ")",
		s.Body.Accept(p).(string),
		")",
	}

	return strings.Join(ret, "")
}

func (p AstPrinter) VisitClassStmt(s internal.ClassStmt) any {
	methods := []string{}
	for _, m := range s.Methods {
		methods = append(methods, p.VisitFuncStmt(m).(string))
	}

	return "(class " + s.Name + "(" + strings.Join(methods, "; ") + ")"
}

func (p AstPrinter) VisitGetExpr(e internal.GetExpr) any {
	return p.parenthesize("call property '"+e.Name+"'", e.Expression)
}
func (p AstPrinter) VisitSetExpr(e internal.SetExpr) any {
	return p.parenthesize("set property", e.Object, e.Value)
}

func (p AstPrinter) VisitReturnStmt(s internal.RreturnStmt) any {
	return p.parenthesize("return", s.Expression)
}

func (p AstPrinter) VisitForSmt(s internal.ForStmt) any {
	ret := []string{"(for"}
	if s.Initializer != nil {
		ret = append(ret, "(initializer", s.Initializer.Accept(p).(string)+")")
	}

	ret = append(ret, "(condition", s.Condition.Accept(p).(string)+")")

	if s.Step != nil {
		ret = append(ret, "(step", s.Step.Accept(p).(string)+")")
	}

	ret = append(ret, "(body", s.Body.Accept(p).(string)+")")
	ret = append(ret, ")")

	return strings.Join(ret, " ")
}

func (p AstPrinter) VisitIfStmt(s internal.IfStmt) any {
	ret := []string{
		p.parenthesize("if", s.Condition).(string),
		s.If.Accept(p).(string),
	}
	if s.Else != nil {
		ret = append(ret, s.Else.Accept(p).(string))
	}

	return strings.Join(ret, " ")
}

func (p AstPrinter) VisitElseStmt(s internal.ElseStmt) any {
	if s.If != nil {
		return p.VisitIfStmt(s.If.(internal.IfStmt))
	} else {
		return "else " + s.Block.Accept(p).(string)
	}
}
func (p AstPrinter) VisitVarStmt(s internal.VarStmt) any {
	return p.parenthesize(s.Name, s.Expression)
}

func (p AstPrinter) VisitPrintStmt(s internal.PrintStmt) any {
	return p.parenthesize("print", s.Expression)
}

func (p AstPrinter) VisitBlockStmt(s internal.BlockStmt) any {
	var ret []string
	for _, s := range s.Stmts {
		ret = append(ret, s.Accept(p).(string))
	}

	return strings.Join(ret, "\n")
}

// TODO: rename to ExpressionStmt
func (p AstPrinter) VisitStmtExpression(s internal.StmtExpression) any {
	return p.parenthesize("stmt", s.Expression)
}

func (p AstPrinter) VisitLogicalExpr(e internal.Logical) any {
	return p.parenthesize(e.Operator.String(), e.Left, e.Right)
}

func (p AstPrinter) VisitAssignmentExpr(e internal.Assignment) any {
	return p.parenthesize(e.Name, e.Expression)
}

func (p AstPrinter) VisitVariableExpr(e internal.Variable) any {
	return e.Name
}

func (p AstPrinter) VisitBinaryExpr(expression internal.Binary) any {
	return p.parenthesize(expression.Operator.String(), expression.Left, expression.Right)
}

func (p AstPrinter) VisitGroupingExpr(expression internal.Grouping) any {
	return p.parenthesize("group", expression.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expression internal.LiteralExpr) any {
	if expression.Value.IsNil() {
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

func (p AstPrinter) VisitUnaryExpr(expression internal.Unary) any {
	return p.parenthesize(expression.Operator.String(), expression.Right)
}

func (p AstPrinter) VisitCallExpr(e internal.Call) any {
	return "call"
}

func (p AstPrinter) parenthesize(name string, expressions ...internal.Expr) any {
	out := bytes.NewBuffer(nil)
	fmt.Fprintf(out, "(%s", name)

	for _, e := range expressions {
		fmt.Fprintf(out, " ")
		fmt.Fprint(out, e.Accept(p))
	}
	fmt.Fprint(out, ")")

	return out.String()
}
