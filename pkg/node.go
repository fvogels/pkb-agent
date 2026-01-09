package pkg

import (
	"pkb-agent/graph/node"

	tea "github.com/charmbracelet/bubbletea"
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

func (node *Node) GetViewer() tea.Model {
	return node.rawNode.GetViewer()
}
