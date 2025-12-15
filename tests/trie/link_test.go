//go:build test

package trie

import (
	"pkb-agent/trie"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLinks(t *testing.T) {
	testCases := [][]string{
		{"a", "b"},
		{"aaa", "aab"},
		{"aab", "aaa"},
		{"aaa", "aab"},
		{"aaa", "aba"},
		{"aaa", "baa"},
		{"aa", "aaa"},
		{"aaa", "aa"},
		{"aa", "aaa", "aaaa"},
		{"aa", "aaa", "aaaab"},
		{"a", "b", "bx", "c"},
		{"a", "b", "cf", "bx", "c"},
		{"a", "ab", "b", "ba"},
		{"b", "ba", "bb", "bc", "bd"},
		{"b", "baaaaa", "bo"},
		{"bash", "bookmark", "bubble"},
	}

	for _, testCase := range testCases {
		t.Run(strings.Join(testCase, "-"), func(t *testing.T) {
			testLinks(t, testCase)
		})
	}
}

func testLinks(t *testing.T, nodes []string) {
	builder := trie.NewBuilder[string]()

	for _, node := range nodes {
		builder.Add(node, node)
	}
	root := builder.Finish()

	slices.Sort(nodes)

	for i := range len(nodes) - 1 {
		n1 := root.Descend(nodes[i])
		n2 := root.Descend(nodes[i+1])

		require.NotNil(t, n1)
		require.NotNil(t, n2)
		require.Equal(t, n2, n1.NextTerminal)
	}

	n := root.Descend(nodes[len(nodes)-1])
	require.NotNil(t, n)
	require.Nil(t, n.NextTerminal)
}
