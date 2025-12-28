//go:build test

package multifile

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	t.Run("a=b", func(t *testing.T) {
		actual, err := parseAttributes("a=b")

		require.NoError(t, err)
		require.Len(t, actual, 1)
		require.Contains(t, actual, "a")
		require.Equal(t, "b", actual["a"])
	})

	t.Run("abc=xyz", func(t *testing.T) {
		actual, err := parseAttributes("abc=xyz")

		require.NoError(t, err)
		require.Len(t, actual, 1)
		require.Contains(t, actual, "abc")
		require.Equal(t, "xyz", actual["abc"])
	})

	t.Run("abc=\"xyz\"", func(t *testing.T) {
		actual, err := parseAttributes(`abc="xyz"`)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		require.Contains(t, actual, "abc")
		require.Equal(t, "xyz", actual["abc"])
	})

	t.Run("abc=\"this is a test\"", func(t *testing.T) {
		actual, err := parseAttributes(`abc="this is a test"`)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		require.Contains(t, actual, "abc")
		require.Equal(t, "this is a test", actual["abc"])
	})

	t.Run("a=b x=y", func(t *testing.T) {
		actual, err := parseAttributes(`a=b x=y`)

		require.NoError(t, err)
		require.Len(t, actual, 2)
		require.Contains(t, actual, "a")
		require.Contains(t, actual, "x")
		require.Equal(t, "b", actual["a"])
		require.Equal(t, "y", actual["x"])
	})
}
