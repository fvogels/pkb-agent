package nodeselectionview

import (
	"pkb-agent/pkg"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgSetRemainingNodes struct {
	RemainingNodes List
	SelectionIndex int
}

type MsgSetSelectedNodes struct {
	SelectedNodes List
}

type msgRemainingNodesWrapper struct {
	wrapped tea.Msg
}

type msgSelectedNodesWrapper struct {
	wrapped tea.Msg
}

type MsgRemainingNodeHighlighted struct {
	Node *pkg.Node
}

type MsgHighlightRemainingNode struct {
	Index int
}
