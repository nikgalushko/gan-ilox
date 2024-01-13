package scanner

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/nikgalushko/gan-ilox/token"
	"github.com/nikgalushko/gan-ilox/token/kind"
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
				token.New(kind.Var, "var", 1, internal.LiteralNil),
				token.New(kind.Identifier, "kek", 1, internal.LiteralNil),
				token.New(kind.Equal, "=", 1, internal.LiteralNil),
				token.New(kind.LeftParen, "(", 1, internal.LiteralNil),
				token.New(kind.Number, "1", 1, internal.NewLiteralInt(1)),
				token.New(kind.Plus, "+", 1, internal.LiteralNil),
				token.New(kind.Number, "2.53", 1, internal.NewLiteralFloat(2.53)),
				token.New(kind.RightParen, ")", 1, internal.LiteralNil),
				token.New(kind.Slash, "/", 1, internal.LiteralNil),
				token.New(kind.Number, "6", 1, internal.NewLiteralInt(6)),
				token.New(kind.Star, "*", 1, internal.LiteralNil),
				token.New(kind.Identifier, "lol", 1, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 1, internal.LiteralNil),
				token.New(kind.EOF, "", 1, internal.LiteralNil),
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
				token.New(kind.Class, "class", 2, internal.LiteralNil),
				token.New(kind.Identifier, "Test", 2, internal.LiteralNil),
				token.New(kind.Less, "<", 2, internal.LiteralNil),
				token.New(kind.Identifier, "Base", 2, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 2, internal.LiteralNil),

				token.New(kind.Identifier, "foo", 3, internal.LiteralNil),
				token.New(kind.LeftParen, "(", 3, internal.LiteralNil),
				token.New(kind.RightParen, ")", 3, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 3, internal.LiteralNil),

				token.New(kind.Super, "super", 4, internal.LiteralNil),
				token.New(kind.Dot, ".", 4, internal.LiteralNil),
				token.New(kind.Identifier, "foo", 4, internal.LiteralNil),
				token.New(kind.LeftParen, "(", 4, internal.LiteralNil),
				token.New(kind.RightParen, ")", 4, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 4, internal.LiteralNil),

				token.New(kind.If, "if", 6, internal.LiteralNil), // if start
				token.New(kind.LeftParen, "(", 6, internal.LiteralNil),
				token.New(kind.This, "this", 6, internal.LiteralNil),
				token.New(kind.Dot, ".", 6, internal.LiteralNil),
				token.New(kind.Identifier, "i", 6, internal.LiteralNil),
				token.New(kind.EqualEqual, "==", 6, internal.LiteralNil),
				token.New(kind.Number, "0", 6, internal.NewLiteralInt(0)),
				token.New(kind.RightParen, ")", 6, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 6, internal.LiteralNil),
				token.New(kind.Return, "return", 7, internal.LiteralNil),
				token.New(kind.Minus, "-", 7, internal.LiteralNil),
				token.New(kind.Number, "1", 7, internal.NewLiteralInt(1)),
				token.New(kind.Semicolon, ";", 7, internal.LiteralNil),
				token.New(kind.RightBrace, "}", 8, internal.LiteralNil), // if end

				token.New(kind.While, "while", 10, internal.LiteralNil), // while start
				token.New(kind.LeftParen, "(", 10, internal.LiteralNil),
				token.New(kind.This, "this", 10, internal.LiteralNil),
				token.New(kind.Dot, ".", 10, internal.LiteralNil),
				token.New(kind.Identifier, "i", 10, internal.LiteralNil),
				token.New(kind.GreaterEqual, ">=", 10, internal.LiteralNil),
				token.New(kind.Number, "100", 10, internal.NewLiteralInt(100)),
				token.New(kind.RightParen, ")", 10, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 10, internal.LiteralNil),
				token.New(kind.This, "this", 11, internal.LiteralNil),
				token.New(kind.Dot, ".", 11, internal.LiteralNil),
				token.New(kind.Identifier, "i", 11, internal.LiteralNil),
				token.New(kind.Equal, "=", 11, internal.LiteralNil),
				token.New(kind.This, "this", 11, internal.LiteralNil),
				token.New(kind.Dot, ".", 11, internal.LiteralNil),
				token.New(kind.Identifier, "i", 11, internal.LiteralNil),
				token.New(kind.Minus, "-", 11, internal.LiteralNil),
				token.New(kind.This, "this", 11, internal.LiteralNil),
				token.New(kind.Dot, ".", 11, internal.LiteralNil),
				token.New(kind.Identifier, "j", 11, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 11, internal.LiteralNil),
				token.New(kind.RightBrace, "}", 12, internal.LiteralNil), // while end

				token.New(kind.For, "for", 13, internal.LiteralNil), // for start
				token.New(kind.LeftParen, "(", 13, internal.LiteralNil),
				token.New(kind.Var, "var", 13, internal.LiteralNil),
				token.New(kind.Identifier, "k", 13, internal.LiteralNil),
				token.New(kind.Equal, "=", 13, internal.LiteralNil),
				token.New(kind.Number, "0", 13, internal.NewLiteralInt(0)),
				token.New(kind.Semicolon, ";", 13, internal.LiteralNil),
				token.New(kind.Identifier, "k", 13, internal.LiteralNil),
				token.New(kind.Less, "<", 13, internal.LiteralNil),
				token.New(kind.This, "this", 13, internal.LiteralNil),
				token.New(kind.Dot, ".", 13, internal.LiteralNil),
				token.New(kind.Identifier, "j", 13, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 13, internal.LiteralNil),
				token.New(kind.Identifier, "k", 13, internal.LiteralNil),
				token.New(kind.Equal, "=", 13, internal.LiteralNil),
				token.New(kind.Identifier, "k", 13, internal.LiteralNil),
				token.New(kind.Plus, "+", 13, internal.LiteralNil),
				token.New(kind.Number, "1", 13, internal.NewLiteralInt(1)),
				token.New(kind.RightParen, ")", 13, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 13, internal.LiteralNil),

				token.New(kind.If, "if", 14, internal.LiteralNil), // if start
				token.New(kind.LeftParen, "(", 14, internal.LiteralNil),
				token.New(kind.Identifier, "k", 14, internal.LiteralNil),
				token.New(kind.Slash, "/", 14, internal.LiteralNil),
				token.New(kind.Number, "2", 14, internal.NewLiteralInt(2)),
				token.New(kind.BangEqual, "!=", 14, internal.LiteralNil),
				token.New(kind.Number, "0", 14, internal.NewLiteralInt(0)),
				token.New(kind.RightParen, ")", 14, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 14, internal.LiteralNil),
				token.New(kind.Print, "print", 15, internal.LiteralNil),
				token.New(kind.Identifier, "k", 15, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 15, internal.LiteralNil),
				token.New(kind.RightBrace, "}", 16, internal.LiteralNil), // if end

				token.New(kind.RightBrace, "}", 17, internal.LiteralNil), // for end

				token.New(kind.RightBrace, "}", 18, internal.LiteralNil), // foo end
				token.New(kind.RightBrace, "}", 19, internal.LiteralNil), // class end

				token.New(kind.EOF, "", 20, internal.LiteralNil),
			},
		},
		{
			in: `if ("str" == "str" and 1 != 1 or k == nil) { print true; } else { print false; }`,
			expected: []token.Token{
				token.New(kind.If, "if", 1, internal.LiteralNil), // if start
				token.New(kind.LeftParen, "(", 1, internal.LiteralNil),
				token.New(kind.String, "str", 1, internal.NewLiteralString("str")),
				token.New(kind.EqualEqual, "==", 1, internal.LiteralNil),
				token.New(kind.String, "str", 1, internal.NewLiteralString("str")),
				token.New(kind.And, "and", 1, internal.LiteralNil),
				token.New(kind.Number, "1", 1, internal.NewLiteralInt(1)),
				token.New(kind.BangEqual, "!=", 1, internal.LiteralNil),
				token.New(kind.Number, "1", 1, internal.NewLiteralInt(1)),
				token.New(kind.Or, "or", 1, internal.LiteralNil),
				token.New(kind.Identifier, "k", 1, internal.LiteralNil),
				token.New(kind.EqualEqual, "==", 1, internal.LiteralNil),
				token.New(kind.Nil, "nil", 1, internal.LiteralNil),
				token.New(kind.RightParen, ")", 1, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 1, internal.LiteralNil), // if body start
				token.New(kind.Print, "print", 1, internal.LiteralNil),
				token.New(kind.True, "true", 1, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 1, internal.LiteralNil),
				token.New(kind.RightBrace, "}", 1, internal.LiteralNil), // if body end
				token.New(kind.Else, "else", 1, internal.LiteralNil),
				token.New(kind.LeftBrace, "{", 1, internal.LiteralNil), // else start
				token.New(kind.Print, "print", 1, internal.LiteralNil),
				token.New(kind.False, "false", 1, internal.LiteralNil),
				token.New(kind.Semicolon, ";", 1, internal.LiteralNil),
				token.New(kind.RightBrace, "}", 1, internal.LiteralNil), // else end
				token.New(kind.EOF, "", 1, internal.LiteralNil),
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
				token.New(kind.Print, "print", 20, internal.LiteralNil),
				token.New(kind.String, "success", 20, internal.NewLiteralString("success")),
				token.New(kind.Semicolon, ";", 20, internal.LiteralNil),
				token.New(kind.EOF, "", 21, internal.LiteralNil),
			},
		},
		{
			in: "var a = 5 & 4 | 3 ^ ~2",
			expected: []token.Token{
				token.New(kind.Var, "var", 1, internal.LiteralNil),
				token.New(kind.Identifier, "a", 1, internal.LiteralNil),
				token.New(kind.Equal, "=", 1, internal.LiteralNil),
				token.New(kind.Number, "5", 1, internal.NewLiteralInt(5)),
				token.New(kind.BitwiseAnd, "&", 1, internal.LiteralNil),
				token.New(kind.Number, "4", 1, internal.NewLiteralInt(4)),
				token.New(kind.BitwiseOr, "|", 1, internal.LiteralNil),
				token.New(kind.Number, "3", 1, internal.NewLiteralInt(3)),
				token.New(kind.BitwiseXor, "^", 1, internal.LiteralNil),
				token.New(kind.BitwiseNot, "~", 1, internal.LiteralNil),
				token.New(kind.Number, "2", 1, internal.NewLiteralInt(2)),
				token.New(kind.EOF, "", 1, internal.LiteralNil),
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
