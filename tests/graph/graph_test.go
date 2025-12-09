//go:build test

package graph

import (
	"pkb-agent/graph"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDescend(t *testing.T) {
	t.Run("Matching", func(t *testing.T) {
		builder := graph.NewBuilder()
		builder.AddNode(&graph.Node{
			Name:  "a",
			Links: nil,
		})
		builder.AddNode(&graph.Node{
			Name:  "b",
			Links: nil,
		})

		g, err := builder.Finish()
		require.NoError(t, err)

		iterator := g.FindMatchingNodes("a")
		require.NotNil(t, iterator.Current())
		require.Equal(t, "a", iterator.Current().Name)
		iterator.Next()
		require.Nil(t, iterator.Current())
	})

	t.Run("Matching 2", func(t *testing.T) {
		builder := graph.NewBuilder()
		builder.AddNode(&graph.Node{
			Name:  "aaa",
			Links: nil,
		})
		builder.AddNode(&graph.Node{
			Name:  "bbb",
			Links: nil,
		})

		g, err := builder.Finish()
		require.NoError(t, err)

		iterator := g.FindMatchingNodes("a")
		require.NotNil(t, iterator.Current())
		require.Equal(t, "aaa", iterator.Current().Name)
		iterator.Next()
		require.Nil(t, iterator.Current())
	})
}
