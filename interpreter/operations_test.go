package interpreter

import (
	"testing"

	"github.com/nikgalushko/gan-ilox/token"
	"github.com/stretchr/testify/require"
)

func TestAdd_OK(t *testing.T) {
	tests := map[string]struct {
		a, b   token.Literal
		expect token.Literal
	}{
		"int + int": {
			a:      token.NewLiteralInt(123),
			b:      token.NewLiteralInt(456),
			expect: token.NewLiteralInt(579),
		},
		"int + float": {
			a:      token.NewLiteralInt(4),
			b:      token.NewLiteralFloat(5.0),
			expect: token.NewLiteralFloat(9.0),
		},
		"float + int": {
			a:      token.NewLiteralFloat(3.2),
			b:      token.NewLiteralInt(5),
			expect: token.NewLiteralFloat(8.2),
		},
		"string + string": {
			a:      token.NewLiteralString("ab"),
			b:      token.NewLiteralString("cd"),
			expect: token.NewLiteralString("abcd"),
		},
	}

	for title, args := range tests {
		t.Run(title, func(t *testing.T) {
			actual, err := add(args.a, args.b)
			require.NoError(t, err)
			require.Equal(t, args.expect, actual)
		})
	}
}

func TestAdd_TypeMissmatch(t *testing.T) {
	possibleArguments := map[string]token.Literal{
		"int":    token.NewLiteralInt(1),
		"float":  token.NewLiteralFloat(1.0),
		"bool":   token.NewLiteralBool(true),
		"string": token.NewLiteralString("1"),
	}
	operations := map[string]func(token.Literal, token.Literal) (token.Literal, error){
		"+": add,
		"-": sub,
		"*": mul,
		"/": div,
	}
	typeMatch := map[string]map[string]struct{}{
		"+": {
			"int+int":       {},
			"int+float":     {},
			"float+int":     {},
			"float+float":   {},
			"string+string": {},
		},
		"-": {
			"int-int":     {},
			"int-float":   {},
			"float-int":   {},
			"float-float": {},
		},
		"*": {
			"int*int":     {},
			"int*float":   {},
			"float*int":   {},
			"float*float": {},
		},
		"/": {
			"int/int":     {},
			"int/float":   {},
			"float/int":   {},
			"float/float": {},
		},
	}

	for operationName, function := range operations {
		for nameA, literalA := range possibleArguments {
			for nameB, litliteralB := range possibleArguments {
				title := nameA + operationName + nameB
				_, err := function(literalA, litliteralB)
				if _, ok := typeMatch[operationName][title]; ok {
					require.NoError(t, err, title)
				} else {
					require.ErrorIs(t, ErrTypeMissmatch, err, title)
				}
			}
		}
	}
}
