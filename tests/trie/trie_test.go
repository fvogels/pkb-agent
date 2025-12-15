//go:build test

package trie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrie(t *testing.T) {
	t.Run("a - b", func(t *testing.T) {
		root := createTrie(
			"a",
			"b",
		)

		nodeA := root.Descend("a")
		nodeB := root.Descend("b")

		require.Equal(t, 1, nodeA.Depth)
		require.Equal(t, nodeB, nodeA.NextTerminal)
		require.Equal(t, 0, nodeA.NextTerminalDepth)

		require.Equal(t, 1, nodeB.Depth)
		require.Nil(t, nodeB.NextTerminal)
	})

	t.Run("a - aa", func(t *testing.T) {
		root := createTrie(
			"a",
			"aa",
		)

		nodeA := root.Descend("a")
		nodeAA := root.Descend("aa")

		require.Equal(t, 1, nodeA.Depth)
		require.Equal(t, nodeAA, nodeA.NextTerminal)
		require.Equal(t, 1, nodeA.NextTerminalDepth)

		require.Equal(t, 2, nodeAA.Depth)
		require.Nil(t, nodeAA.NextTerminal)
	})

	t.Run("a - ba", func(t *testing.T) {
		root := createTrie(
			"a",
			"ba",
		)

		nodeA := root.Descend("a")
		nodeBA := root.Descend("ba")

		require.Equal(t, 1, nodeA.Depth)
		require.Equal(t, nodeBA, nodeA.NextTerminal)
		require.Equal(t, 0, nodeA.NextTerminalDepth)

		require.Equal(t, 2, nodeBA.Depth)
		require.Nil(t, nodeBA.NextTerminal)
	})

	t.Run("a - ba - bb", func(t *testing.T) {
		root := createTrie(
			"a",
			"ba",
			"bb",
		)

		nodeA := root.Descend("a")
		nodeBA := root.Descend("ba")
		nodeBB := root.Descend("bb")

		require.Equal(t, 1, nodeA.Depth)
		require.Equal(t, nodeBA, nodeA.NextTerminal)
		require.Equal(t, 0, nodeA.NextTerminalDepth)

		require.Equal(t, 2, nodeBA.Depth)
		require.Equal(t, nodeBA, nodeA.NextTerminal)
		require.Equal(t, 1, nodeBA.NextTerminalDepth)

		require.Equal(t, 2, nodeBB.Depth)
		require.Nil(t, nodeBB.NextTerminal)
	})

	t.Run("ab - b", func(t *testing.T) {
		root := createTrie(
			"ab",
			"b",
		)

		nodeAB := root.Descend("ab")
		nodeB := root.Descend("b")

		require.Equal(t, 2, nodeAB.Depth)
		require.Equal(t, nodeB, nodeAB.NextTerminal)
		require.Equal(t, 0, nodeAB.NextTerminalDepth)

		require.Equal(t, 1, nodeB.Depth)
		require.Nil(t, nodeB.NextTerminal)
	})
}
