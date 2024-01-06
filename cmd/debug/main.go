package main

import (
	"fmt"

	"github.com/nikgalushko/gan-ilox/debug"
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

	fmt.Println(debug.AstPrinter{E: e})
}
