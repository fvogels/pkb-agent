package model

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/util/queue"
	"pkb-agent/util/set"
	"strings"
)

// determineIntersectionNodes computes which nodes are compatible with the selected nodes and the search filter.
func determineIntersectionNodes(input string, graph *pkg.Graph, selectedNodes list.List[*pkg.Node], includeAncestors bool, includeIndirectDescendants bool, onlyLeaves bool) list.List[*pkg.Node] {
	intersectionNodes := collectDescendantIntersection(graph, selectedNodes, includeIndirectDescendants)

	if includeAncestors {
		intersectionNodes = collectAncestorsUnion(graph, intersectionNodes)
	}

	searchStrings := strings.Split(input, " ")
	intersectionNodesAncestors := collectNodesMatchingSearches(graph, searchStrings)
	intersectionNodes.IntersectWith(intersectionNodesAncestors)

	list.ForEach(selectedNodes, func(index int, selectedNode *pkg.Node) {
		intersectionNodes.Remove(selectedNode.GetIndex())
	})

	result := []*pkg.Node{}
	for i := range graph.GetNodeCount() {
		if intersectionNodes.Contains(i) {
			node := graph.FindNodeByIndex(i)

			if !onlyLeaves || len(node.GetBacklinks()) == 0 {
				result = append(result, node)
			}
		}
	}

	return list.FromSlice(result)
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
	intersection := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())
	for i := range graph.GetNodeCount() {
		intersection.Add(i)
	}

	union := set.NewIntSetWithInitialCapacity(graph.GetNodeCount())

	for i := range graph.GetNodeCount() {
		if roots.Contains(i) {
			ancestors := collectAncestors(graph, graph.FindNodeByIndex(i))
			union.UnionWith(ancestors)
			intersection.IntersectWith(ancestors)
		}
	}

	union.DifferenceWith(intersection)
	return union
}
