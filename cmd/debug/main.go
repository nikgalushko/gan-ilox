package main

import (
	"bytes"
	"fmt"

	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

func main() {
	e := expr.Binary[string]{
		Left: expr.Unary[string]{
			Operator: token.New(token.Minus, "-", 1, nil),
			Right:    expr.Literal[string]{Value: 123},
		},
		Operator: token.New(token.Star, "*", 1, nil),
		Right: expr.Grouping[string]{
			Expression: expr.Literal[string]{Value: 45.67},
		},
	}

	fmt.Println(AstPrinter{e: e})
}

type AstPrinter struct {
	e expr.Expr[string]
}

func (p AstPrinter) String() string {
	return p.e.Accept(p)
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
