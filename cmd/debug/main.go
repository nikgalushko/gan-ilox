package main

import (
	"fmt"

	"github.com/nikgalushko/gan-ilox/debug"
	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/interpreter"
	"github.com/nikgalushko/gan-ilox/token"
)

func main() {
	e := expr.Binary{
		Left: expr.Unary{
			Operator: token.New(token.Minus, "-", 1, nil),
			Right:    expr.Literal{Value: 123},
		},
		Operator: token.New(token.Star, "*", 1, nil),
		Right: expr.Grouping{
			Expression: expr.Literal{Value: 45},
		},
	}

	fmt.Println(debug.AstPrinter{E: e})
	fmt.Println(interpreter.New(e).Eval())
}
