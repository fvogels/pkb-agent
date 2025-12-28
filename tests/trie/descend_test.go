//go:build test

package trie

import (
	"pkb-agent/util/trie"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDescend(t *testing.T) {
	t.Run("Single terminal", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("abc", "abc")
		root := builder.Finish()

		node := root.Descend("abc")

		require.Len(t, node.Terminals, 1)
		require.Equal(t, "abc", node.Terminals[0])
		require.Nil(t, node.NextTerminal)
	})

	t.Run("Prefix", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("abc", "abc")
		root := builder.Finish()

		node := root.Descend("ab")

		require.Len(t, node.Terminals, 0)
		require.Equal(t, root.Descend("abc"), node.NextTerminal)
	})

	t.Run("Nonexistent prefix", func(t *testing.T) {
		builder := trie.NewBuilder[string]()
		builder.Add("abc", "abc")
		root := builder.Finish()

		node := root.Descend("x")

		require.Nil(t, node)
	})
}
