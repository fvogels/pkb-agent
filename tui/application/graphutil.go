package application

import (
	"pkb-agent/graph"
	"pkb-agent/util"
	"pkb-agent/util/set"
	"slices"
	"sort"
	"strings"
)

// determineIntersectionNodes computes which nodes are compatible with the selected nodes and the search filter.
func determineIntersectionNodes(input string, g *graph.Graph, selectedNodes []*graph.Node, includeLinked bool, includeIndirectAncestors bool) []*graph.Node {
	if len(selectedNodes) == 0 {
		// Deal with this case separately for efficiency reasons
		iterator := g.FindMatchingNodes(input)

		// nameSet is used to prevent duplicates
		// Adding the selected nodes ensures that already selected nodes do not appear as remaining choices
		nameSet := set.FromSlice(util.Map(selectedNodes, func(node *graph.Node) string { return node.GetName() }))
		remaining := []*graph.Node{}

		for iterator.Current() != nil {
			// The same node can occur more than once during iteration
			// Ensure that we add each node only once to remainingNodes
			name := iterator.Current().GetName()
			if nameSet.Contains(name) {
				iterator.Next()
				continue
			}

			nameSet.Add(name)
			remaining = append(remaining, iterator.Current())
			iterator.Next()
		}

		sort.Slice(remaining, func(i, j int) bool {
			return strings.ToLower(remaining[i].GetName()) < strings.ToLower(remaining[j].GetName())
		})

		return remaining
	}

	// Step 1: collect the intersection of the descendants of the selected nodes
	remainingNodeSet := collectDescendants(g, selectedNodes[0], includeIndirectAncestors)

	for _, selectedNode := range selectedNodes[1:] {
		remainingNodeSet.Intersect(collectDescendants(g, selectedNode, includeIndirectAncestors))
	}

	// Step 2: optionally collect all ancestors of the remaining nodes
	// Note: currently, only direct ancestors are selected
	// This could be improved to also include indirect ancestors
	if includeLinked {
		for _, node := range remainingNodeSet.ToSlice() {
			for _, linkedNode := range g.FindNodeByIndex(node).GetLinks() {
				linkedNodeIndex := linkedNode.GetIndex()
				if !remainingNodeSet.Contains(linkedNodeIndex) {
					remainingNodeSet.Add(linkedNodeIndex)
				}
			}
		}
	}

	// Step 3: only keep nodes that are compatible with the filter
	if len(input) != 0 {
		subselection := set.New[int]()

		iterator := g.FindMatchingNodes(input)

		for iterator.Current() != nil {
			if remainingNodeSet.Contains(iterator.Current().GetIndex()) {
				subselection.Add(iterator.Current().GetIndex())
			}
			iterator.Next()
		}

		remainingNodeSet = subselection
	}

	// Step 4: remove nodes that are already selected
	for _, selectedNode := range selectedNodes {
		remainingNodeSet.Remove(selectedNode.GetIndex())
	}

	result := remainingNodeSet.ToSlice()
	slices.Sort(result)

	return util.Map(result, func(index int) *graph.Node {
		return g.FindNodeByIndex(index)
	})
}

// collectDescendants collects the names of all backlinked nodes
func collectDescendants(g *graph.Graph, node *graph.Node, includeIndirect bool) set.Set[int] {
	result := set.New[int]()
	queue := make([]*graph.Node, 1, 20)
	queue[0] = node

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for _, backlinked := range current.GetBacklinks() {
			backlinkedIndex := backlinked.GetIndex()
			result.Add(backlinkedIndex)

			if includeIndirect {
				descendant := g.FindNodeByIndex(backlinkedIndex)
				queue = append(queue, descendant)
			}
		}
	}

	return result
}
