package debug

import (
	"bytes"
	"fmt"

	"github.com/nikgalushko/gan-ilox/expr"
)

type AstPrinter struct {
	E expr.Expr
}

func (p AstPrinter) String() string {
	return p.E.Accept(p).(string)
}

func (p AstPrinter) VisitBinaryExpr(expression expr.Binary) any {
	return p.parenthesize(expression.Operator.Lexeme, expression.Left, expression.Right)
}

func (p AstPrinter) VisitGroupingExpr(expression expr.Grouping) any {
	return p.parenthesize("group", expression.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expression expr.Literal) any {
	if expression.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", expression.Value)
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
