package tabp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParser(t *testing.T) {
	t.Run("Number", func(t *testing.T) {
		t.Run("Float", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("3.14"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.Equal(t, 3.14, v)
		})

		t.Run("Integer", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("1000"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.Equal(t, 1000, v)
		})

	})

	t.Run("Symbol", func(t *testing.T) {
		t.Run("PI", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("PI"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.Equal(t, Symbol("PI"), v)
		})

		t.Run("WithSpaces", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("|A SYMBOL WITH SPACES|"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.Equal(t, Symbol("|A SYMBOL WITH SPACES|"), v)
		})
	})

	t.Run("String", func(t *testing.T) {
		parser := NewParser(bytes.NewBufferString(`"foo bar baz"`))

		v, err := parser.Parse()
		require.NoError(t, err.Cause)
		require.Equal(t, "foo bar baz", v)
	})

	t.Run("Table", func(t *testing.T) {
		t.Run("SingleSymbol", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("(Symbol)"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.IsType(t, &Table{}, v)
			require.Equal(t, Symbol("SYMBOL"), v.(*Table).Get(0))
		})

		t.Run("SingleNumber", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString("(3.14)"))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.IsType(t, &Table{}, v)
			require.Equal(t, 3.14, v.(*Table).Get(0))
		})

		t.Run("SingleString", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString(`("my string")`))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.IsType(t, &Table{}, v)
			require.Equal(t, "my string", v.(*Table).Get(0))
		})

		t.Run("Mixed", func(t *testing.T) {
			parser := NewParser(bytes.NewBufferString(`(Symbol 3.14 "my string" foo: "bar")`))

			v, err := parser.Parse()
			require.NoError(t, err.Cause)
			require.IsType(t, &Table{}, v)
			require.Equal(t, Symbol("SYMBOL"), v.(*Table).Get(0))
			require.Equal(t, 3.14, v.(*Table).Get(1))
			require.Equal(t, "my string", v.(*Table).Get(2))
			require.Equal(t, "bar", v.(*Table).Get(Symbol("FOO")))
		})
	})
}
