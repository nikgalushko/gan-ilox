package main

import (
	"fmt"

	"github.com/nikgalushko/gan-ilox/debug"
	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token/kind"
)

func main() {
	e := internal.Binary{
		Left: internal.Unary{
			Operator: kind.Minus,
			Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(123)},
		},
		Operator: kind.Star,
		Right: internal.Grouping{
			Expression: internal.Binary{
				Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(45)},
				Operator: kind.Plus,
				Right:    internal.Variable{Name: "kek"},
			},
		},
	}

	fmt.Println(debug.AstPrinter{E: e})
}
