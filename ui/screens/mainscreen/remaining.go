package mainscreen

import (
	"pkb-agent/graph"
	"pkb-agent/util"
	"pkb-agent/util/set"
	"slices"
	"sort"
	"strings"
)

func determineRemainingNodes(input string, g *graph.Graph, selectedNodes []*graph.Node, includeLinked bool) []*graph.Node {
	iterator := g.FindMatchingNodes(input)

	// nameSet is used to prevent duplicates
	// Adding the selected nodes ensures that already selected nodes do not appear as remaining choices
	nameSet := set.FromSlice(util.Map(selectedNodes, func(node *graph.Node) string { return node.Name }))
	remaining := []*graph.Node{}

	for iterator.Current() != nil {
		// The same node can occur more than once during iteration
		// Ensure that we add each node only once to remainingNodes
		name := iterator.Current().Name
		if nameSet.Contains(name) {
			iterator.Next()
			continue
		}

		if !util.All(selectedNodes, func(selectedNode *graph.Node) bool {
			return slices.Contains(iterator.Current().Links, selectedNode.Name)
		}) {
			iterator.Next()
			continue
		}

		nameSet.Add(name)
		remaining = append(remaining, iterator.Current())
		iterator.Next()
	}

	if includeLinked {
		imax := len(remaining)
		for i := 0; i != imax; i++ {
			node := remaining[i]

			for _, linkedNodeName := range node.Links {
				if !nameSet.Contains(linkedNodeName) {
					nameSet.Add(linkedNodeName)
					linkedNode := g.FindNode(linkedNodeName)
					remaining = append(remaining, linkedNode)
				}
			}
		}
	}

	sort.Slice(remaining, func(i, j int) bool {
		return strings.ToLower(remaining[i].Name) < strings.ToLower(remaining[j].Name)
	})

	return remaining
}
