package graph

import (
	"pkb-agent/trie"
)

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
		if iterator.current != nil && iterator.current.NextTerminalDepth < iterator.minimalDepth {
			iterator.current = nil
		} else {
			iterator.current = iterator.current.NextTerminal
			iterator.index = 0
		}
	}
}
