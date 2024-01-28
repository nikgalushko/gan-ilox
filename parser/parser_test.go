package parser

import (
	"testing"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/scanner"
	"github.com/nikgalushko/gan-ilox/token/kind"
	"github.com/stretchr/testify/require"
)

func TestParser_Skip(t *testing.T) {
	const code = `
	var a = 1;
	a - return
	for;
	a=2;
	`
	tokens, err := scanner.NewScanner(code).ScanTokens()
	require.NoError(t, err)

	stmts, err := New(tokens).Parse()

	require.Equal(t, "expect expression", err.(PraseError)[0].Error())
	require.Equal(t, "expect '(' after for", err.(PraseError)[1].Error())
	require.Equal(t, []internal.Stmt{
		internal.VarStmt{
			Name:       "a",
			Expression: internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
		},
		internal.StmtExpression{
			Expression: internal.Assignment{
				Name:       "a",
				Expression: internal.LiteralExpr{Value: internal.NewLiteralInt(2)},
			},
		},
	}, stmts)
}

func TestParser_HappyPath(t *testing.T) {
	for _, args := range testCases {
		t.Run(args.Name, func(t *testing.T) {
			tokens, err := scanner.NewScanner(args.Code).ScanTokens()
			require.NoError(t, err)

			stmts, err := New(tokens).Parse()
			if args.Err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, args.ExpectedStmt, stmts)
			}
		})
	}
}

type Case struct {
	Name         string
	Code         string
	ExpectedStmt []internal.Stmt
	Err          bool
}

var (
	testCases = []Case{
		classDeclaration,
		varDeclaration,
		functionDeclaration,
		returnStatement,
		reqturnStatement2,
		ifStatement,
		forStatement,
		forStatement2,
		forStatement3,
		assignment,
		assignment2,
		logical,
		equality,
		unary,
		functionCall,
	}
	classDeclaration = Case{
		Name: "class declaration",
		Code: `
			class Foo {
				method0() {
					print "method0";
				}

				method1(a) {
					print a;
				}

				method2(a, b) {
					print a + b;
				}
			}
		`,
		ExpectedStmt: []internal.Stmt{
			internal.ClassStmt{
				Name: "Foo",
				Methods: []internal.FuncStmt{
					{
						Name:       "method0",
						Parameters: nil,
						Body: internal.BlockStmt{
							Stmts: []internal.Stmt{
								internal.PrintStmt{
									Expression: internal.LiteralExpr{
										Value: internal.NewLiteralString("method0"),
									},
								},
							},
						},
					},
					{
						Name:       "method1",
						Parameters: []string{"a"},
						Body: internal.BlockStmt{
							Stmts: []internal.Stmt{
								internal.PrintStmt{
									Expression: internal.Variable{Name: "a"},
								},
							},
						},
					},
					{
						Name:       "method2",
						Parameters: []string{"a", "b"},
						Body: internal.BlockStmt{
							Stmts: []internal.Stmt{
								internal.PrintStmt{
									Expression: internal.Binary{
										Left:     internal.Variable{Name: "a"},
										Operator: kind.Plus,
										Right:    internal.Variable{Name: "b"},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	varDeclaration = Case{
		Name: "var declaration",
		Code: `
		var b;
		var a = 5-e;
		`,
		ExpectedStmt: []internal.Stmt{
			internal.VarStmt{Name: "b"},
			internal.VarStmt{
				Name: "a",
				Expression: internal.Binary{
					Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(5)},
					Operator: kind.Minus,
					Right:    internal.Variable{Name: "e"},
				},
			},
		},
	}
	functionDeclaration = Case{
		Name: "function declaration",
		Code: `
		fun foo() {}
		fun bar(a) {}
		fun baz(a,b,c,d) {}
		`,
		ExpectedStmt: []internal.Stmt{
			internal.FuncStmt{
				Name: "foo",
				Body: internal.BlockStmt{},
			},
			internal.FuncStmt{
				Name:       "bar",
				Parameters: []string{"a"},
				Body:       internal.BlockStmt{},
			},
			internal.FuncStmt{
				Name:       "baz",
				Parameters: []string{"a", "b", "c", "d"},
				Body:       internal.BlockStmt{},
			},
		},
	}
	returnStatement = Case{
		Name: "return statement outsied a function",
		Code: `return 1;`,
		Err:  true,
	}
	reqturnStatement2 = Case{
		Name: "return statement insede a function",
		Code: `fun foo() { return 1; }`,
		ExpectedStmt: []internal.Stmt{
			internal.FuncStmt{
				Name: "foo",
				Body: internal.BlockStmt{
					Stmts: []internal.Stmt{
						internal.RreturnStmt{
							Expression: internal.LiteralExpr{
								Value: internal.NewLiteralInt(1),
							},
						},
					},
				},
			},
		},
	}
	ifStatement = Case{
		Name: "if else if else",
		Code: `if (1==1) {a;} else if (1>1) {b;} else {c;}`,
		ExpectedStmt: []internal.Stmt{
			internal.IfStmt{
				Condition: internal.Binary{
					Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
					Operator: kind.EqualEqual,
					Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
				},
				If: internal.BlockStmt{
					Stmts: []internal.Stmt{internal.StmtExpression{Expression: internal.Variable{Name: "a"}}},
				},
				Else: internal.IfStmt{
					Condition: internal.Binary{
						Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
						Operator: kind.Greater,
						Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
					},
					If: internal.BlockStmt{
						Stmts: []internal.Stmt{internal.StmtExpression{Expression: internal.Variable{Name: "b"}}},
					},
					Else: internal.BlockStmt{
						Stmts: []internal.Stmt{internal.StmtExpression{Expression: internal.Variable{Name: "c"}}},
					},
				},
			},
		},
	}
	forStatement = Case{
		Name: "classic for statement",
		Code: `for (var a = 1; a < 10; a = a + 1) { print a; }`,
		ExpectedStmt: []internal.Stmt{
			internal.ForStmt{
				Initializer: internal.VarStmt{Name: "a", Expression: internal.LiteralExpr{Value: internal.NewLiteralInt(1)}},
				Condition: internal.Binary{
					Left:     internal.Variable{Name: "a"},
					Operator: kind.Less,
					Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(10)},
				},
				Step: internal.Assignment{
					Name: "a",
					Expression: internal.Binary{
						Left:     internal.Variable{Name: "a"},
						Operator: kind.Plus,
						Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
					},
				},
				Body: internal.BlockStmt{
					Stmts: []internal.Stmt{
						internal.PrintStmt{
							Expression: internal.Variable{Name: "a"},
						},
					},
				},
			},
		},
	}
	forStatement2 = Case{
		Name: "for-loop only condition",
		Code: `for (i < 10) {}`,
		ExpectedStmt: []internal.Stmt{
			internal.ForStmt{
				Condition: internal.Binary{
					Left:     internal.Variable{Name: "i"},
					Operator: kind.Less,
					Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(10)},
				},
				Body: internal.BlockStmt{},
			},
		},
	}
	forStatement3 = Case{
		Name: "infinite loop",
		Code: `for {}`,
		ExpectedStmt: []internal.Stmt{
			internal.ForStmt{
				Body: internal.BlockStmt{},
			},
		},
	}
	assignment = Case{
		Name: "simple assignment",
		Code: `i = 5;`,
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.Assignment{Name: "i", Expression: internal.LiteralExpr{Value: internal.NewLiteralInt(5)}},
			},
		},
	}
	assignment2 = Case{
		Name: "set field to object",
		Code: `foo.kek = 5;`,
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.SetExpr{
					Name:   "kek",
					Object: internal.Variable{Name: "foo"},
					Value:  internal.LiteralExpr{Value: internal.NewLiteralInt(5)},
				},
			},
		},
	}
	logical = Case{
		Name: "logical",
		Code: `(1 and 5) or (2 and 6);`,
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.Logical{
					Left: internal.Grouping{
						Expression: internal.Logical{
							Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
							Operator: kind.And,
							Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(5)},
						},
					},
					Operator: kind.Or,
					Right: internal.Grouping{
						Expression: internal.Logical{
							Left:     internal.LiteralExpr{Value: internal.NewLiteralInt(2)},
							Operator: kind.And,
							Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(6)},
						},
					},
				},
			},
		},
	}
	equality = Case{
		Name: "equality",
		Code: `a == b != c;`,
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.Binary{
					Left: internal.Binary{
						Left:     internal.Variable{Name: "a"},
						Operator: kind.EqualEqual,
						Right:    internal.Variable{Name: "b"},
					},
					Operator: kind.BangEqual,
					Right:    internal.Variable{Name: "c"},
				},
			},
		},
	}
	unary = Case{
		Name: "unary",
		Code: `!a;-1;~1;`,
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.Unary{
					Operator: kind.Bang,
					Right:    internal.Variable{Name: "a"},
				},
			},
			internal.StmtExpression{
				Expression: internal.Unary{
					Operator: kind.Minus,
					Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
				},
			},
			internal.StmtExpression{
				Expression: internal.Unary{
					Operator: kind.BitwiseNot,
					Right:    internal.LiteralExpr{Value: internal.NewLiteralInt(1)},
				},
			},
		},
	}
	functionCall = Case{
		Name: "function call",
		Code: "foo(a);",
		ExpectedStmt: []internal.Stmt{
			internal.StmtExpression{
				Expression: internal.Call{
					Arguments: []internal.Expr{internal.Variable{Name: "a"}},
					Callee:    internal.Variable{Name: "foo"},
				},
			},
		},
	}
)
