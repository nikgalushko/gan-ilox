package env

import (
	"testing"

	"github.com/nikgalushko/gan-ilox/internal"
	"github.com/stretchr/testify/require"
)

func TestEnvironment(t *testing.T) {
	e := New()

	v, err := e.Get("test")
	require.ErrorIs(t, ErrUndefinedVariable, err)
	require.Equal(t, internal.LiteralNil, v)

	err = e.Assign("test", internal.NewLiteralInt(123))
	require.ErrorIs(t, ErrUndefinedVariable, err)

	e.Define("test", internal.NewLiteralString("str"))

	v, err = e.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsString(), "str")

	err = e.Assign("test", internal.NewLiteralInt(123))
	require.NoError(t, err)

	v, err = e.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsInt(), int64(123))

	e2 := NewWithParent(e)

	v, err = e2.Get("test")
	require.NoError(t, err)
	require.Equal(t, v.AsInt(), int64(123))

	err = e2.Assign("test", internal.NewLiteralBool(true))
	require.NoError(t, err)

	v, err = e2.Get("test")
	require.NoError(t, err)
	require.True(t, v.AsBool())

	v, err = e.Get("test")
	require.NoError(t, err)
	require.True(t, v.AsBool())

	e2.Define("kek", internal.NewLiteralString("e2_lol"))
	v, err = e2.Get("kek")
	require.NoError(t, err)
	require.Equal(t, "e2_lol", v.AsString())

	v, err = e.Get("kek")
	require.ErrorIs(t, ErrUndefinedVariable, err)
	require.Equal(t, internal.LiteralNil, v)
}
