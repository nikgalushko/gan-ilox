package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nikgalushko/gan-ilox/token"
)

func TestScanTokens(t *testing.T) {
	tests := []struct {
		in       string
		expected []token.Token
		err      error
	}{
		{
			in: "var kek = (1 + 2.53)/6 * lol;",
			expected: []token.Token{
				token.New(token.Var, "var", 1, token.LiteralNil),
				token.New(token.Identifier, "kek", 1, token.LiteralNil),
				token.New(token.Equal, "=", 1, token.LiteralNil),
				token.New(token.LeftParen, "(", 1, token.LiteralNil),
				token.New(token.Number, "1", 1, token.NewLiteralInt(1)),
				token.New(token.Plus, "+", 1, token.LiteralNil),
				token.New(token.Number, "2.53", 1, token.NewLiteralFloat(2.53)),
				token.New(token.RightParen, ")", 1, token.LiteralNil),
				token.New(token.Slash, "/", 1, token.LiteralNil),
				token.New(token.Number, "6", 1, token.NewLiteralInt(6)),
				token.New(token.Star, "*", 1, token.LiteralNil),
				token.New(token.Identifier, "lol", 1, token.LiteralNil),
				token.New(token.Semicolon, ";", 1, token.LiteralNil),
				token.New(token.EOF, "", 1, token.LiteralNil),
			},
		},
		{
			in: `// my test class
				class Test < Base {
					foo() {
						super.foo();

						if (this.i == 0) {
							return -1;
						}

						while (this.i >= 100) {
							this.i = this.i - this.j;
						}
						for (var k = 0; k < this.j; k = k + 1) {
							if (k / 2 != 0) {
								print k;
							}
						}
					}
				}
			`,
			expected: []token.Token{
				token.New(token.Class, "class", 2, token.LiteralNil),
				token.New(token.Identifier, "Test", 2, token.LiteralNil),
				token.New(token.Less, "<", 2, token.LiteralNil),
				token.New(token.Identifier, "Base", 2, token.LiteralNil),
				token.New(token.LeftBrace, "{", 2, token.LiteralNil),

				token.New(token.Identifier, "foo", 3, token.LiteralNil),
				token.New(token.LeftParen, "(", 3, token.LiteralNil),
				token.New(token.RightParen, ")", 3, token.LiteralNil),
				token.New(token.LeftBrace, "{", 3, token.LiteralNil),

				token.New(token.Super, "super", 4, token.LiteralNil),
				token.New(token.Dot, ".", 4, token.LiteralNil),
				token.New(token.Identifier, "foo", 4, token.LiteralNil),
				token.New(token.LeftParen, "(", 4, token.LiteralNil),
				token.New(token.RightParen, ")", 4, token.LiteralNil),
				token.New(token.Semicolon, ";", 4, token.LiteralNil),

				token.New(token.If, "if", 6, token.LiteralNil), // if start
				token.New(token.LeftParen, "(", 6, token.LiteralNil),
				token.New(token.This, "this", 6, token.LiteralNil),
				token.New(token.Dot, ".", 6, token.LiteralNil),
				token.New(token.Identifier, "i", 6, token.LiteralNil),
				token.New(token.EqualEqual, "==", 6, token.LiteralNil),
				token.New(token.Number, "0", 6, token.NewLiteralInt(0)),
				token.New(token.RightParen, ")", 6, token.LiteralNil),
				token.New(token.LeftBrace, "{", 6, token.LiteralNil),
				token.New(token.Return, "return", 7, token.LiteralNil),
				token.New(token.Minus, "-", 7, token.LiteralNil),
				token.New(token.Number, "1", 7, token.NewLiteralInt(1)),
				token.New(token.Semicolon, ";", 7, token.LiteralNil),
				token.New(token.RightBrace, "}", 8, token.LiteralNil), // if end

				token.New(token.While, "while", 10, token.LiteralNil), // while start
				token.New(token.LeftParen, "(", 10, token.LiteralNil),
				token.New(token.This, "this", 10, token.LiteralNil),
				token.New(token.Dot, ".", 10, token.LiteralNil),
				token.New(token.Identifier, "i", 10, token.LiteralNil),
				token.New(token.GreaterEqual, ">=", 10, token.LiteralNil),
				token.New(token.Number, "100", 10, token.NewLiteralInt(100)),
				token.New(token.RightParen, ")", 10, token.LiteralNil),
				token.New(token.LeftBrace, "{", 10, token.LiteralNil),
				token.New(token.This, "this", 11, token.LiteralNil),
				token.New(token.Dot, ".", 11, token.LiteralNil),
				token.New(token.Identifier, "i", 11, token.LiteralNil),
				token.New(token.Equal, "=", 11, token.LiteralNil),
				token.New(token.This, "this", 11, token.LiteralNil),
				token.New(token.Dot, ".", 11, token.LiteralNil),
				token.New(token.Identifier, "i", 11, token.LiteralNil),
				token.New(token.Minus, "-", 11, token.LiteralNil),
				token.New(token.This, "this", 11, token.LiteralNil),
				token.New(token.Dot, ".", 11, token.LiteralNil),
				token.New(token.Identifier, "j", 11, token.LiteralNil),
				token.New(token.Semicolon, ";", 11, token.LiteralNil),
				token.New(token.RightBrace, "}", 12, token.LiteralNil), // while end

				token.New(token.For, "for", 13, token.LiteralNil), // for start
				token.New(token.LeftParen, "(", 13, token.LiteralNil),
				token.New(token.Var, "var", 13, token.LiteralNil),
				token.New(token.Identifier, "k", 13, token.LiteralNil),
				token.New(token.Equal, "=", 13, token.LiteralNil),
				token.New(token.Number, "0", 13, token.NewLiteralInt(0)),
				token.New(token.Semicolon, ";", 13, token.LiteralNil),
				token.New(token.Identifier, "k", 13, token.LiteralNil),
				token.New(token.Less, "<", 13, token.LiteralNil),
				token.New(token.This, "this", 13, token.LiteralNil),
				token.New(token.Dot, ".", 13, token.LiteralNil),
				token.New(token.Identifier, "j", 13, token.LiteralNil),
				token.New(token.Semicolon, ";", 13, token.LiteralNil),
				token.New(token.Identifier, "k", 13, token.LiteralNil),
				token.New(token.Equal, "=", 13, token.LiteralNil),
				token.New(token.Identifier, "k", 13, token.LiteralNil),
				token.New(token.Plus, "+", 13, token.LiteralNil),
				token.New(token.Number, "1", 13, token.NewLiteralInt(1)),
				token.New(token.RightParen, ")", 13, token.LiteralNil),
				token.New(token.LeftBrace, "{", 13, token.LiteralNil),

				token.New(token.If, "if", 14, token.LiteralNil), // if start
				token.New(token.LeftParen, "(", 14, token.LiteralNil),
				token.New(token.Identifier, "k", 14, token.LiteralNil),
				token.New(token.Slash, "/", 14, token.LiteralNil),
				token.New(token.Number, "2", 14, token.NewLiteralInt(2)),
				token.New(token.BangEqual, "!=", 14, token.LiteralNil),
				token.New(token.Number, "0", 14, token.NewLiteralInt(0)),
				token.New(token.RightParen, ")", 14, token.LiteralNil),
				token.New(token.LeftBrace, "{", 14, token.LiteralNil),
				token.New(token.Print, "print", 15, token.LiteralNil),
				token.New(token.Identifier, "k", 15, token.LiteralNil),
				token.New(token.Semicolon, ";", 15, token.LiteralNil),
				token.New(token.RightBrace, "}", 16, token.LiteralNil), // if end

				token.New(token.RightBrace, "}", 17, token.LiteralNil), // for end

				token.New(token.RightBrace, "}", 18, token.LiteralNil), // foo end
				token.New(token.RightBrace, "}", 19, token.LiteralNil), // class end

				token.New(token.EOF, "", 20, token.LiteralNil),
			},
		},
		{
			in: `if ("str" == "str" and 1 != 1 or k == nil) { print true; } else { print false; }`,
			expected: []token.Token{
				token.New(token.If, "if", 1, token.LiteralNil), // if start
				token.New(token.LeftParen, "(", 1, token.LiteralNil),
				token.New(token.String, "str", 1, token.NewLiteralString("str")),
				token.New(token.EqualEqual, "==", 1, token.LiteralNil),
				token.New(token.String, "str", 1, token.NewLiteralString("str")),
				token.New(token.And, "and", 1, token.LiteralNil),
				token.New(token.Number, "1", 1, token.NewLiteralInt(1)),
				token.New(token.BangEqual, "!=", 1, token.LiteralNil),
				token.New(token.Number, "1", 1, token.NewLiteralInt(1)),
				token.New(token.Or, "or", 1, token.LiteralNil),
				token.New(token.Identifier, "k", 1, token.LiteralNil),
				token.New(token.EqualEqual, "==", 1, token.LiteralNil),
				token.New(token.Nil, "nil", 1, token.LiteralNil),
				token.New(token.RightParen, ")", 1, token.LiteralNil),
				token.New(token.LeftBrace, "{", 1, token.LiteralNil), // if body start
				token.New(token.Print, "print", 1, token.LiteralNil),
				token.New(token.True, "true", 1, token.LiteralNil),
				token.New(token.Semicolon, ";", 1, token.LiteralNil),
				token.New(token.RightBrace, "}", 1, token.LiteralNil), // if body end
				token.New(token.Else, "else", 1, token.LiteralNil),
				token.New(token.LeftBrace, "{", 1, token.LiteralNil), // else start
				token.New(token.Print, "print", 1, token.LiteralNil),
				token.New(token.False, "false", 1, token.LiteralNil),
				token.New(token.Semicolon, ";", 1, token.LiteralNil),
				token.New(token.RightBrace, "}", 1, token.LiteralNil), // else end
				token.New(token.EOF, "", 1, token.LiteralNil),
			},
		},
		{
			in: `// my test class
			/*class Test < Base {
				foo() {
					super.foo();

					if (this.i == 0) {
						return -1;
					}
					//inner comment
					while (this.i >= 100) {
						this.i = this.i - this.j;
					}
					for (var k = 0; k < this.j; k = k + 1) {
						if (k / 2 != 0) {
							print k;
						}
					}
				}
			}*/
			print "success";
			`,
			expected: []token.Token{
				token.New(token.Print, "print", 20, token.LiteralNil),
				token.New(token.String, "success", 20, token.NewLiteralString("success")),
				token.New(token.Semicolon, ";", 20, token.LiteralNil),
				token.New(token.EOF, "", 21, token.LiteralNil),
			},
		},
		{
			in: "var a = 5 & 4 | 3 ^ ~2",
			expected: []token.Token{
				token.New(token.Var, "var", 1, token.LiteralNil),
				token.New(token.Identifier, "a", 1, token.LiteralNil),
				token.New(token.Equal, "=", 1, token.LiteralNil),
				token.New(token.Number, "5", 1, token.NewLiteralInt(5)),
				token.New(token.BitwiseAnd, "&", 1, token.LiteralNil),
				token.New(token.Number, "4", 1, token.NewLiteralInt(4)),
				token.New(token.BitwiseOr, "|", 1, token.LiteralNil),
				token.New(token.Number, "3", 1, token.NewLiteralInt(3)),
				token.New(token.BitwiseXor, "^", 1, token.LiteralNil),
				token.New(token.BitwiseNot, "~", 1, token.LiteralNil),
				token.New(token.Number, "2", 1, token.NewLiteralInt(2)),
				token.New(token.EOF, "", 1, token.LiteralNil),
			},
		},
	}

	for _, args := range tests {
		s := NewScanner(args.in)
		actually, err := s.ScanTokens()

		if args.err != nil {
			require.Error(t, err)
		} else {
			require.Equal(t, args.expected, actually)
		}
	}
}
