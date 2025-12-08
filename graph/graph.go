package graph

import (
	"maps"
	"pkb-agent/trie"
	"pkb-agent/util"
)

type Graph struct {
	nodes    map[string]*Node
	trieRoot *trie.Node[*Node]
}

func (graph *Graph) FindNode(name string) *Node {
	node, ok := graph.nodes[name]

	if !ok {
		return nil
	}

	return node
}

func (graph *Graph) ListNodeNames() []string {
	result := []string{}
	generator := maps.Keys(graph.nodes)
	generator(util.CollectTo(&result))

	return result
}
