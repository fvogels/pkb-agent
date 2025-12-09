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

func (graph *Graph) FindNameMatches(str string) MatchIterator {
	trieNode := graph.trieRoot.Descend(str)

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
		minimalDepth: len(str),
	}
}

type MatchIterator struct {
	minimalDepth int
	current      *trie.Node[*Node]
	index        int
}

func (iterator *MatchIterator) Current() *Node {
	if iterator.current == nil {
		return nil
	}

	return iterator.current.Terminals[iterator.index]
}

func (iterator *MatchIterator) Next() {
	iterator.index++

	if iterator.index == len(iterator.current.Terminals) {
		iterator.current = iterator.current.NextTerminal
		iterator.index = 0

		if iterator.current != nil && iterator.current.NextTerminalDepth <= iterator.minimalDepth {
			iterator.current = nil
		}
	}
}
