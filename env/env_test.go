package env

import (
	"testing"

	"github.com/nikgalushko/gan-ilox/token"
	"github.com/stretchr/testify/require"
)

func TestEnvironment(t *testing.T) {
	e := New()

	v, err := e.Get("test")
	require.ErrorIs(t, ErrUndefinedVariable, err)
	require.Equal(t, token.LiteralNil, v)

	err = e.Assign("test", token.NewLiteralInt(123))
	require.ErrorIs(t, ErrUndefinedVariable, err)

	e.Define("test", token.NewLiteralString("str"))

	v, err = e.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsString(), "str")

	err = e.Assign("test", token.NewLiteralInt(123))
	require.NoError(t, err)

	v, err = e.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsInt(), int64(123))

	e2 := NewWithParent(e)

	v, err = e2.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsInt(), int64(123))

	err = e2.Assign("test", token.NewLiteralBool(true))
	require.NoError(t, err)

	v, err = e2.Get("test")
	require.NoError(t, err)
	require.True(t, v.AsBool())

	v, err = e.Get("test")
	require.NoError(t, err)
	require.True(t, v.AsBool())

	e2.Define("kek", token.NewLiteralString("e2_lol"))
	v, err = e2.Get("kek")
	require.NoError(t, err)
	require.Equal(t, "e2_lol", v.AsString())

	v, err = e.Get("kek")
	require.ErrorIs(t, ErrUndefinedVariable, err)
	require.Equal(t, token.LiteralNil, v)
}
