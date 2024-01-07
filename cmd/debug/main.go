package main

import (
	"fmt"

	"github.com/nikgalushko/gan-ilox/debug"
	"github.com/nikgalushko/gan-ilox/expr"
	"github.com/nikgalushko/gan-ilox/token"
)

func main() {
	e := expr.Binary{
		Left: expr.Unary{
			Operator: token.New(token.Minus, "-", 1, token.LiteralNil),
			Right:    expr.Literal{Value: token.NewLiteralInt(123)},
		},
		Operator: token.New(token.Star, "*", 1, token.LiteralNil),
		Right: expr.Grouping{
			Expression: expr.Binary{
				Left:     expr.Literal{Value: token.NewLiteralInt(45)},
				Operator: token.New(token.Plus, "+", 1, token.LiteralNil),
				Right:    expr.Variable{Name: token.New(token.Identifier, "kek", 1, token.LiteralNil)},
			},
		},
	}

	fmt.Println(debug.AstPrinter{E: e})
}
