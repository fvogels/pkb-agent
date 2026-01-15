package pkg

import (
	"pkb-agent/pkg/node"
	"pkb-agent/tui"
)

type Node struct {
	id        int
	links     []*Node
	backlinks []*Node
	rawNode   node.RawNode
}

func (node *Node) GetName() string {
	return node.rawNode.GetName()
}

func (node *Node) GetLinks() []*Node {
	return node.links
}

func (node *Node) GetBacklinks() []*Node {
	return node.backlinks
}

func (node *Node) GetIndex() int {
	return node.id
}

func (node *Node) GetViewer(messageQueue tui.MessageQueue) tui.Component {
	return node.rawNode.CreateViewer(messageQueue)
}
