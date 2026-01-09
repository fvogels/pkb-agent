package pkg

import "pkb-agent/util/set"

func (graph *Graph) FindRedundantLinks(node *Node) set.Set[int] {
	indirectAncestors := set.New[int]()
	directAncestors := set.New[int]()

	for _, directAncestor := range node.links {
		directAncestors.Add(directAncestor.id)

		graph.CollectAncestors(directAncestor, func(indirectAncestor *Node) {
			indirectAncestors.Add(indirectAncestor.id)
		})
	}

	directAncestors.Intersect(indirectAncestors)

	return directAncestors
}
