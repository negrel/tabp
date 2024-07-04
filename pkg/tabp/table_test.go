package tabp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTable(t *testing.T) {
	t.Run("GetSet", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			tab := Table{}

			tab.Set("foo", "bar")
			require.Equal(t, "bar", tab.Get("foo"))
		})

		t.Run("UpdateValue", func(t *testing.T) {
			tab := Table{}

			tab.Set("foo", "bar")
			tab.Set("foo", "baz")
			require.Equal(t, "baz", tab.Get("foo"))
		})
	})

	t.Run("Append", func(t *testing.T) {
		t.Run("Empty", func(t *testing.T) {
			tab := Table{}

			size := tab.Append("foo")
			require.Equal(t, 1, size)

			require.Equal(t, "foo", tab.Get(0))
		})
		t.Run("AfterSequence", func(t *testing.T) {
			tab := Table{}

			tab.Set(0, "foo")

			size := tab.Append("bar")
			require.Equal(t, 2, size)

			require.Equal(t, "bar", tab.Get(1))
		})

		t.Run("FillHoleInSequence", func(t *testing.T) {
			tab := Table{}

			tab.Set(0, "foo")
			tab.Set(2, "baz")

			size := tab.Append("bar")
			require.Equal(t, 3, size)

			require.Equal(t, "bar", tab.Get(1))
		})
	})

	t.Run("SequenceLength", func(t *testing.T) {
		tab := Table{}

		for i := 0; i < 1_000_000; i++ {
			tab.Set(i, i)
		}

		require.Equal(t, 1_000_000, tab.SequenceLen())
	})
}
