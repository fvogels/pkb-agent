package pkg

import (
	"maps"
	"pkb-agent/util"
	"pkb-agent/util/trie"

	_ "pkb-agent/pkg/loaders/glob"
	_ "pkb-agent/pkg/loaders/sequence"
	_ "pkb-agent/pkg/loaders/tagger"
	_ "pkb-agent/pkg/node/atom"
	_ "pkb-agent/pkg/node/backblaze"
	_ "pkb-agent/pkg/node/bookmark"
	_ "pkb-agent/pkg/node/hybrid"
)

type Graph struct {
	nodesByIndex []*Node
	nodesByName  map[string]*Node
	trieRoot     *trie.Node[*Node]
}

func (graph *Graph) FindNodeByIndex(index int) *Node {
	return graph.nodesByIndex[index]
}

func (graph *Graph) FindNodeByName(name string) *Node {
	node, ok := graph.nodesByName[name]

	if !ok {
		return nil
	}

	return node
}

func (graph *Graph) GetNodeCount() int {
	return len(graph.nodesByName)
}

func (graph *Graph) ListNodeNames() []string {
	result := []string{}
	generator := maps.Keys(graph.nodesByName)
	generator(util.CollectTo(&result))

	return result
}

func (graph *Graph) ListNodes() []*Node {
	nodes := make([]*Node, len(graph.nodesByIndex))
	copy(nodes, graph.nodesByIndex)
	return nodes
}

func (graph *Graph) FindMatchingNodes(nameMatch string) MatchIterator {
	trieNode := graph.trieRoot.Descend(nameMatch)

	if trieNode == nil {
		return MatchIterator{
			current: nil,
		}
	}

	if len(trieNode.Terminals) == 0 {
		next := trieNode.NextTerminal

		if next.Depth >= trieNode.Depth {
			trieNode = next
		}
	}

	return MatchIterator{
		current:      trieNode,
		index:        0,
		minimalDepth: len(nameMatch),
	}
}

func (graph *Graph) CollectAncestors(node *Node, yield func(*Node)) {
	stackCapacity := 10
	stack := make([]*Node, 1, stackCapacity)
	stack[0] = node

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1 : stackCapacity]

		for _, ancestor := range current.links {
			yield(ancestor)
			stack = append(stack, ancestor)
		}
	}
}
