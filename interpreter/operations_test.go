package interpreter

import (
	"math/rand"
	"testing"

	"github.com/nikgalushko/gan-ilox/token"
	"github.com/stretchr/testify/require"
)

func TestArithmeticNumbers(t *testing.T) {
	operations := map[string]func(token.Literal, token.Literal) (token.Literal, error){
		"+": add,
		"-": sub,
		"*": mul,
		"/": div,
	}

	goldenInts := map[string]func(a, b int64) int64{
		"+": func(a, b int64) int64 { return a + b },
		"-": func(a, b int64) int64 { return a - b },
		"*": func(a, b int64) int64 { return a * b },
		"/": func(a, b int64) int64 { return a / b },
	}
	goldenFloats := map[string]func(a, b float64) float64{
		"+": func(a, b float64) float64 { return a + b },
		"-": func(a, b float64) float64 { return a - b },
		"*": func(a, b float64) float64 { return a * b },
		"/": func(a, b float64) float64 { return a / b },
	}
	possibleArguments := []string{"int", "float"}

	newLiteralOfType := func(_type string) token.Literal {
		switch _type {
		case "int":
			return token.NewLiteralInt(rand.Int63())
		case "float":
			return token.NewLiteralFloat(rand.Float64())
		}
		t.Fatal("undefined type " + _type)
		return token.LiteralNil
	}

	for opName, op := range operations {
		for _, aType := range possibleArguments {
			for _, bType := range possibleArguments {
				a := newLiteralOfType(aType)
				b := newLiteralOfType(bType)
				actual, err := op(a, b)
				require.NoError(t, err)

				if aType == bType && aType == "int" {
					require.Equal(t, goldenInts[opName](a.AsInt(), b.AsInt()), actual.AsInt())
				} else {
					require.InDelta(t, goldenFloats[opName](a.AsFloat(), b.AsFloat()), actual.AsFloat(), 1e-9)
				}
			}
		}
	}
}
func TestTypeMissmatch(t *testing.T) {
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
