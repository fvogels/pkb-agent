package model

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/util"
	"pkb-agent/util/queue"
	"pkb-agent/util/set"
	"slices"
	"sort"
	"strings"
)

// determineIntersectionNodes computes which nodes are compatible with the selected nodes and the search filter.
func determineIntersectionNodes(input string, graph *pkg.Graph, selectedNodes list.List[*pkg.Node], includeLinked bool, includeIndirectAncestors bool) list.List[*pkg.Node] {
	r1 := collectDescendantIntersection(graph, selectedNodes, includeIndirectAncestors)

	if includeLinked {
		r1 = collectAncestorsUnion(graph, r1)
	}

	searchStrings := strings.Split(input, " ")
	r2 := collectNodesMatchingSearches(graph, searchStrings)
	r1.IntersectWith(r2)

	result := []*pkg.Node{}
	for i := range graph.GetNodeCount() {
		if r1.Contains(i) {
			node := graph.FindNodeByIndex(i)
			result = append(result, node)
		}
	}

	return list.FromSlice(result)
}

func determineIntersectionNodesOld(input string, graph *pkg.Graph, selectedNodes list.List[*pkg.Node], includeLinked bool, includeIndirectAncestors bool) list.List[*pkg.Node] {
	if selectedNodes.Size() == 0 {
		// Deal with this case separately for efficiency reasons
		iterator := graph.FindMatchingNodes(input)

		// nameSet is used to prevent duplicates
		// Adding the selected nodes ensures that already selected nodes do not appear as remaining choices
		nameSet := set.FromSlice(util.Map(list.ToSlice(selectedNodes), func(node *pkg.Node) string { return node.GetName() }))
		remaining := []*pkg.Node{}

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

		return list.FromSlice(remaining)
	}

	// Step 1: collect the intersection of the descendants of the selected nodes
	remainingNodeSet := collectDescendants(graph, selectedNodes.At(0), includeIndirectAncestors)

	for i := 1; i < selectedNodes.Size(); i++ {
		selectedNode := selectedNodes.At(i)
		remainingNodeSet.Intersect(collectDescendants(graph, selectedNode, includeIndirectAncestors))
	}

	// Step 2: optionally collect all ancestors of the remaining nodes
	// Note: currently, only direct ancestors are selected
	// This could be improved to also include indirect ancestors
	if includeLinked {
		for _, node := range remainingNodeSet.ToSlice() {
			for _, linkedNode := range graph.FindNodeByIndex(node).GetLinks() {
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

		iterator := graph.FindMatchingNodes(input)

		for iterator.Current() != nil {
			if remainingNodeSet.Contains(iterator.Current().GetIndex()) {
				subselection.Add(iterator.Current().GetIndex())
			}
			iterator.Next()
		}

		remainingNodeSet = subselection
	}

	// Step 4: remove nodes that are already selected
	list.ForEach(selectedNodes, func(_ int, selectedNode *pkg.Node) {
		remainingNodeSet.Remove(selectedNode.GetIndex())
	})

	intersectionNodeIndices := remainingNodeSet.ToSlice()
	slices.Sort(intersectionNodeIndices)

	intersectionNodes := util.Map(intersectionNodeIndices, func(index int) *pkg.Node {
		return graph.FindNodeByIndex(index)
	})

	return list.FromSlice(intersectionNodes)
}

// collectDescendants collects the names of all backlinked nodes
func collectDescendants(graph *pkg.Graph, node *pkg.Node, includeIndirect bool) set.Set[int] {
	result := set.New[int]()
	queue := make([]*pkg.Node, 1, 20)
	queue[0] = node

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for _, backlinked := range current.GetBacklinks() {
			backlinkedIndex := backlinked.GetIndex()
			result.Add(backlinkedIndex)

			if includeIndirect {
				descendant := graph.FindNodeByIndex(backlinkedIndex)
				queue = append(queue, descendant)
			}
		}
	}

	return result
}

// collectDescendants2 collects the names of all backlinked nodes
func collectDescendants2(graph *pkg.Graph, node *pkg.Node, includeIndirect bool) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())
	queue := make([]*pkg.Node, 1, 20)
	queue[0] = node

	for len(queue) > 0 {
		current := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		for _, backlinked := range current.GetBacklinks() {
			backlinkedIndex := backlinked.GetIndex()
			result.Add(backlinkedIndex)

			if includeIndirect {
				descendant := graph.FindNodeByIndex(backlinkedIndex)
				queue = append(queue, descendant)
			}
		}
	}

	return result
}

func collectDescendantIntersection(graph *pkg.Graph, nodes list.List[*pkg.Node], includeIndirect bool) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())

	for i := range graph.GetNodeCount() {
		result.Add(i)
	}

	for i := range nodes.Size() {
		node := nodes.At(i)
		r := collectDescendants2(graph, node, includeIndirect)
		result.IntersectWith(r)
	}

	return result
}

func collectNodesMatchingSearch(graph *pkg.Graph, searchString string) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())
	iterator := graph.FindMatchingNodes(searchString)

	for iterator.Current() != nil {
		result.Add(iterator.Current().GetIndex())
		iterator.Next()
	}

	return result
}

func collectNodesMatchingSearches(graph *pkg.Graph, searchStrings []string) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())

	for i := range graph.GetNodeCount() {
		result.Add(i)
	}

	for _, searchString := range searchStrings {
		r := collectNodesMatchingSearch(graph, searchString)
		result.IntersectWith(r)
	}

	return result
}

func collectAncestors(graph *pkg.Graph, node *pkg.Node) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())
	queue := queue.New[int]()
	queue.Enqueue(node.GetIndex())

	for !queue.IsEmpty() {
		nodeID := queue.Dequeue()
		if !result.Contains(nodeID) {
			result.Add(nodeID)

			node := graph.FindNodeByIndex(nodeID)
			for _, parent := range node.GetLinks() {
				queue.Enqueue(parent.GetIndex())
			}
		}
	}

	return result
}

func collectAncestorsUnion(graph *pkg.Graph, roots *set.IntSet) *set.IntSet {
	result := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())

	for i := range graph.GetNodeCount() {
		if roots.Contains(i) {
			r := collectAncestors(graph, graph.FindNodeByIndex(i))
			result.UnionWith(r)
		}
	}

	return result
}
