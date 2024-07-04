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

			tab.Set(Symbol("foo"), "bar")
			tab.Set(Symbol("foo"), "baz")
			require.Equal(t, "baz", tab.Get(Symbol("foo")))

			require.Equal(t, `(foo: "baz")`, Sexpr(&tab))
		})

		t.Run("DeleteElementInSequence", func(t *testing.T) {
			tab := Table{}

			tab.Set(0, 0)
			tab.Set(1, 1)
			tab.Set(2, 2)
			require.Equal(t, 3, tab.SeqLen())

			// Delete element.
			tab.Set(1, nil)

			require.Equal(t, 0, tab.Get(0))
			require.Equal(t, nil, tab.Get(1))
			require.Equal(t, 2, tab.Get(2))

			require.Equal(t, `(0 2: 2)`, Sexpr(&tab))
			require.Equal(t, 1, tab.SeqLen())
		})

		t.Run("FillHoleInSequence", func(t *testing.T) {
			tab := Table{}

			tab.Set(4, 4)
			tab.Set(5, 5)
			tab.Set(2, 2)
			tab.Set(1, 1)
			tab.Set(0, 0)
			require.Equal(t, 3, tab.SeqLen())

			// Fill hole in sequence.
			tab.Set(3, 3)

			require.Equal(t, 6, tab.SeqLen())
			require.Equal(t, 6, len(tab.array)) // Elements are stored in array.
			// Key are present because of insertion order.
			require.Equal(t, `(4: 4 5: 5 2 1: 1 0: 0 3: 3)`, Sexpr(&tab))
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

		require.Equal(t, 1_000_000, tab.SeqLen())
	})
}

func BenchmarkTable(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tab := Table{}

			for i := 0; i < 10000; i++ {
				tab.Append(i)
			}
		}
	})

}

func BenchmarkSlice(b *testing.B) {
	b.Run("Append", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var tab []int

			for i := 0; i < 10000; i++ {
				tab = append(tab, i)
			}
		}
	})
}
