package debug

import (
	"bytes"
	"fmt"

	"github.com/nikgalushko/gan-ilox/expr"
)

type AstPrinter struct {
	E expr.Expr[string]
}

func (p AstPrinter) String() string {
	return p.E.Accept(p)
}

func (p AstPrinter) VisitBinaryExpr(expression expr.Binary[string]) string {
	return p.parenthesize(expression.Operator.Lexeme, expression.Left, expression.Right)
}

func (p AstPrinter) VisitGroupingExpr(expression expr.Grouping[string]) string {
	return p.parenthesize("group", expression.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expression expr.Literal[string]) string {
	if expression.Value == nil {
		return "nil"
	}

	return fmt.Sprintf("%v", expression.Value)
}

func (p AstPrinter) VisitUnaryExpr(expression expr.Unary[string]) string {
	return p.parenthesize(expression.Operator.Lexeme, expression.Right)
}

func (p AstPrinter) parenthesize(name string, expressions ...expr.Expr[string]) string {
	out := bytes.NewBuffer(nil)
	fmt.Fprintf(out, "(%s", name)

	for _, e := range expressions {
		fmt.Fprintf(out, " ")
		fmt.Fprint(out, e.Accept(p))
	}
	fmt.Fprint(out, ")")

	return out.String()
}
