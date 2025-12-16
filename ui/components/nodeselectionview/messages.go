package nodeselectionview

import (
	"pkb-agent/graph"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgSetRemainingNodes struct {
	RemainingNodes List
	SelectionIndex int
}

type MsgSetSelectedNodes struct {
	SelectedNodes List
}

type MsgSelectNext struct{}

type MsgSelectPrevious struct{}

type msgRemainingNodesWrapper struct {
	wrapped tea.Msg
}

type msgSelectedNodesWrapper struct {
	wrapped tea.Msg
}

type MsgRemainingNodeHighlighted struct {
	Node *graph.Node
}

type MsgSelectRemainingNode struct {
	Index int
}
