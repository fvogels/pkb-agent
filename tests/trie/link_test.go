//go:build test

package trie

import (
	"pkb-agent/trie"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinks(t *testing.T) {
	t.Run("aaa - aab", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("aaa", "aaa")
		builder.Add("aab", "aab")
		root := builder.Finish()

		node1 := root.Descend("aaa")
		node2 := root.Descend("aab")

		require.NotNil(t, node1)
		require.NotNil(t, node2)
		require.Equal(t, node2, node1.NextTerminal)
		require.Equal(t, 3, node1.NextTerminalDepth)
		require.Nil(t, node2.NextTerminal)
	})

	t.Run("aaa - aba", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("aaa", "aaa")
		builder.Add("aba", "aba")
		root := builder.Finish()

		node1 := root.Descend("aaa")
		node2 := root.Descend("aba")

		require.NotNil(t, node1)
		require.NotNil(t, node2)
		require.Equal(t, node2, node1.NextTerminal)
		require.Equal(t, 2, node1.NextTerminalDepth)
		require.Nil(t, node2.NextTerminal)
	})

	t.Run("aaa - baa", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("aaa", "aaa")
		builder.Add("baa", "baa")
		root := builder.Finish()

		node1 := root.Descend("aaa")
		node2 := root.Descend("baa")

		require.NotNil(t, node1)
		require.NotNil(t, node2)
		require.Equal(t, node2, node1.NextTerminal)
		require.Equal(t, 1, node1.NextTerminalDepth)
		require.Nil(t, node2.NextTerminal)
	})
}
