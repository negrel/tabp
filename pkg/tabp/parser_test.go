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
	})
}
