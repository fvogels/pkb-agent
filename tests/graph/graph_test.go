//go:build test

package graph

import (
	"pkb-agent/graph"
	"testing"

	"pkb-agent/tests/testlib"

	"github.com/stretchr/testify/require"
)

func TestDescend(t *testing.T) {
	t.Run("Matching", func(t *testing.T) {
		builder := graph.NewBuilder()
		builder.AddNode(&testlib.TestNode{
			Name:     "a",
			Keywords: []string{"a"},
		})
		builder.AddNode(&testlib.TestNode{
			Name:     "b",
			Keywords: []string{"b"},
		})

		g, err := builder.Finish()
		require.NoError(t, err)

		iterator := g.FindMatchingNodes("a")
		require.NotNil(t, iterator.Current())
		require.Equal(t, "a", iterator.Current().GetName())
		iterator.Next()
		require.Nil(t, iterator.Current())
	})

	t.Run("Matching 2", func(t *testing.T) {
		builder := graph.NewBuilder()
		builder.AddNode(&testlib.TestNode{
			Name:     "aaa",
			Keywords: []string{"aaa"},
		})
		builder.AddNode(&testlib.TestNode{
			Name:     "bbb",
			Keywords: []string{"abb"},
		})

		g, err := builder.Finish()
		require.NoError(t, err)

		iterator := g.FindMatchingNodes("a")
		require.NotNil(t, iterator.Current())
		require.Equal(t, "aaa", iterator.Current().GetName())
		iterator.Next()
		require.Nil(t, iterator.Current())
	})
}
