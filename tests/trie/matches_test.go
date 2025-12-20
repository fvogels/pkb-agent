//go:build test

package trie

import (
	"fmt"
	"pkb-agent/graph"
	"pkb-agent/graph/nodes/atom"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type matchTestCase struct {
	nodes        []string
	searchString string
	expected     []string
}

func TestGraphSearch(t *testing.T) {
	testCases := []matchTestCase{
		{
			nodes: []string{
				"a",
				"b",
				"c",
			},
			searchString: "",
			expected: []string{
				"a",
				"b",
				"c",
			},
		},
		{
			nodes: []string{
				"a",
				"b",
				"c",
			},
			searchString: "a",
			expected: []string{
				"a",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.string(), func(t *testing.T) {
			testCase.test(t)
		})
	}
}

func (tc matchTestCase) test(t *testing.T) {
	builder := graph.NewBuilder()

	for _, node := range tc.nodes {
		builder.AddNode(&graph.Node{
			Name:  node,
			Links: nil,
			Type:  "atom",
			Info:  atom.Info{},
		})
	}

	g, err := builder.Finish()
	require.NoError(t, err)

	iterator := g.FindMatchingNodes(tc.searchString)

	for index, current := range tc.expected {
		require.NotNil(t, iterator.Current(), "index %d, expecting %s", index, current)
		require.Equal(t, current, iterator.Current().Name)
		iterator.Next()
	}

	require.Nil(t, iterator.Current())
}

func (tc matchTestCase) string() string {
	return fmt.Sprintf(
		"Nodes: %s - Search: \"%s\"",
		strings.Join(tc.nodes, ", "),
		tc.searchString,
	)
}
