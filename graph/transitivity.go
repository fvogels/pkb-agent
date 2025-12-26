package graph

import "pkb-agent/util/set"

func (graph *Graph) FindRedundantLinks(node *Node) set.Set[int] {
	indirectAncestors := set.New[int]()
	directAncestors := set.New[int]()

	for _, directAncestorName := range node.Links {
		directAncestor := graph.FindNodeByName(directAncestorName)
		directAncestors.Add(directAncestor.Index)

		graph.CollectAncestors(directAncestor, func(indirectAncestor *Node) {
			indirectAncestors.Add(indirectAncestor.Index)
		})
	}

	directAncestors.Intersect(indirectAncestors)

	return directAncestors
}
