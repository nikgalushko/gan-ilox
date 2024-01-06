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
				token.New(token.Var, "var", 1, nil),
				token.New(token.Identifier, "kek", 1, nil),
				token.New(token.Equal, "=", 1, nil),
				token.New(token.LeftParen, "(", 1, nil),
				token.New(token.Number, "1", 1, float64(1)),
				token.New(token.Plus, "+", 1, nil),
				token.New(token.Number, "2.53", 1, float64(2.53)),
				token.New(token.RightParen, ")", 1, nil),
				token.New(token.Slash, "/", 1, nil),
				token.New(token.Number, "6", 1, float64(6)),
				token.New(token.Star, "*", 1, nil),
				token.New(token.Identifier, "lol", 1, nil),
				token.New(token.Semicolon, ";", 1, nil),
				token.New(token.EOF, "", 1, nil),
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
				token.New(token.Class, "class", 2, nil),
				token.New(token.Identifier, "Test", 2, nil),
				token.New(token.Less, "<", 2, nil),
				token.New(token.Identifier, "Base", 2, nil),
				token.New(token.LeftBrace, "{", 2, nil),

				token.New(token.Identifier, "foo", 3, nil),
				token.New(token.LeftParen, "(", 3, nil),
				token.New(token.RightParen, ")", 3, nil),
				token.New(token.LeftBrace, "{", 3, nil),

				token.New(token.Super, "super", 4, nil),
				token.New(token.Dot, ".", 4, nil),
				token.New(token.Identifier, "foo", 4, nil),
				token.New(token.LeftParen, "(", 4, nil),
				token.New(token.RightParen, ")", 4, nil),
				token.New(token.Semicolon, ";", 4, nil),

				token.New(token.If, "if", 6, nil), // if start
				token.New(token.LeftParen, "(", 6, nil),
				token.New(token.This, "this", 6, nil),
				token.New(token.Dot, ".", 6, nil),
				token.New(token.Identifier, "i", 6, nil),
				token.New(token.EqualEqual, "==", 6, nil),
				token.New(token.Number, "0", 6, float64(0)),
				token.New(token.RightParen, ")", 6, nil),
				token.New(token.LeftBrace, "{", 6, nil),
				token.New(token.Return, "return", 7, nil),
				token.New(token.Minus, "-", 7, nil),
				token.New(token.Number, "1", 7, float64(1)),
				token.New(token.Semicolon, ";", 7, nil),
				token.New(token.RightBrace, "}", 8, nil), // if end

				token.New(token.While, "while", 10, nil), // while start
				token.New(token.LeftParen, "(", 10, nil),
				token.New(token.This, "this", 10, nil),
				token.New(token.Dot, ".", 10, nil),
				token.New(token.Identifier, "i", 10, nil),
				token.New(token.GreaterEqual, ">=", 10, nil),
				token.New(token.Number, "100", 10, float64(100)),
				token.New(token.RightParen, ")", 10, nil),
				token.New(token.LeftBrace, "{", 10, nil),
				token.New(token.This, "this", 11, nil),
				token.New(token.Dot, ".", 11, nil),
				token.New(token.Identifier, "i", 11, nil),
				token.New(token.Equal, "=", 11, nil),
				token.New(token.This, "this", 11, nil),
				token.New(token.Dot, ".", 11, nil),
				token.New(token.Identifier, "i", 11, nil),
				token.New(token.Minus, "-", 11, nil),
				token.New(token.This, "this", 11, nil),
				token.New(token.Dot, ".", 11, nil),
				token.New(token.Identifier, "j", 11, nil),
				token.New(token.Semicolon, ";", 11, nil),
				token.New(token.RightBrace, "}", 12, nil), // while end

				token.New(token.For, "for", 13, nil), // for start
				token.New(token.LeftParen, "(", 13, nil),
				token.New(token.Var, "var", 13, nil),
				token.New(token.Identifier, "k", 13, nil),
				token.New(token.Equal, "=", 13, nil),
				token.New(token.Number, "0", 13, float64(0)),
				token.New(token.Semicolon, ";", 13, nil),
				token.New(token.Identifier, "k", 13, nil),
				token.New(token.Less, "<", 13, nil),
				token.New(token.This, "this", 13, nil),
				token.New(token.Dot, ".", 13, nil),
				token.New(token.Identifier, "j", 13, nil),
				token.New(token.Semicolon, ";", 13, nil),
				token.New(token.Identifier, "k", 13, nil),
				token.New(token.Equal, "=", 13, nil),
				token.New(token.Identifier, "k", 13, nil),
				token.New(token.Plus, "+", 13, nil),
				token.New(token.Number, "1", 13, float64(1)),
				token.New(token.RightParen, ")", 13, nil),
				token.New(token.LeftBrace, "{", 13, nil),

				token.New(token.If, "if", 14, nil), // if start
				token.New(token.LeftParen, "(", 14, nil),
				token.New(token.Identifier, "k", 14, nil),
				token.New(token.Slash, "/", 14, nil),
				token.New(token.Number, "2", 14, float64(2)),
				token.New(token.BangEqual, "!=", 14, nil),
				token.New(token.Number, "0", 14, float64(0)),
				token.New(token.RightParen, ")", 14, nil),
				token.New(token.LeftBrace, "{", 14, nil),
				token.New(token.Print, "print", 15, nil),
				token.New(token.Identifier, "k", 15, nil),
				token.New(token.Semicolon, ";", 15, nil),
				token.New(token.RightBrace, "}", 16, nil), // if end

				token.New(token.RightBrace, "}", 17, nil), // for end

				token.New(token.RightBrace, "}", 18, nil), // foo end
				token.New(token.RightBrace, "}", 19, nil), // class end

				token.New(token.EOF, "", 20, nil),
			},
		},
		{
			in: `if ("str" == "str" and 1 != 1 or k == nil) { print true; } else { print false; }`,
			expected: []token.Token{
				token.New(token.If, "if", 1, nil), // if start
				token.New(token.LeftParen, "(", 1, nil),
				token.New(token.String, "str", 1, "str"),
				token.New(token.EqualEqual, "==", 1, nil),
				token.New(token.String, "str", 1, "str"),
				token.New(token.And, "and", 1, nil),
				token.New(token.Number, "1", 1, float64(1)),
				token.New(token.BangEqual, "!=", 1, nil),
				token.New(token.Number, "1", 1, float64(1)),
				token.New(token.Or, "or", 1, nil),
				token.New(token.Identifier, "k", 1, nil),
				token.New(token.EqualEqual, "==", 1, nil),
				token.New(token.Nil, "nil", 1, nil),
				token.New(token.RightParen, ")", 1, nil),
				token.New(token.LeftBrace, "{", 1, nil), // if body start
				token.New(token.Print, "print", 1, nil),
				token.New(token.True, "true", 1, nil),
				token.New(token.Semicolon, ";", 1, nil),
				token.New(token.RightBrace, "}", 1, nil), // if body end
				token.New(token.Else, "else", 1, nil),
				token.New(token.LeftBrace, "{", 1, nil), // else start
				token.New(token.Print, "print", 1, nil),
				token.New(token.False, "false", 1, nil),
				token.New(token.Semicolon, ";", 1, nil),
				token.New(token.RightBrace, "}", 1, nil), // else end
				token.New(token.EOF, "", 1, nil),
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
				token.New(token.Print, "print", 20, nil),
				token.New(token.String, "success", 20, "success"),
				token.New(token.Semicolon, ";", 20, nil),
				token.New(token.EOF, "", 21, nil),
			},
		},
		{
			in: "var a = 5 & 4 | 3 ^ ~2",
			expected: []token.Token{
				token.New(token.Var, "var", 1, nil),
				token.New(token.Identifier, "a", 1, nil),
				token.New(token.Equal, "=", 1, nil),
				token.New(token.Number, "5", 1, float64(5)),
				token.New(token.BitwiseAnd, "&", 1, nil),
				token.New(token.Number, "4", 1, float64(4)),
				token.New(token.BitwiseOr, "|", 1, nil),
				token.New(token.Number, "3", 1, float64(3)),
				token.New(token.BitwiseXor, "^", 1, nil),
				token.New(token.BitwiseNot, "~", 1, nil),
				token.New(token.Number, "2", 1, float64(2)),
				token.New(token.EOF, "", 1, nil),
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
