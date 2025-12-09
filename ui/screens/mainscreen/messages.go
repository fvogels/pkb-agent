package mainscreen

import (
	"pkb-agent/graph"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgGraphLoaded struct {
	graph *graph.Graph
}

type MsgUpdateNodeList struct{}

type msgToSelectableNodeView struct {
	wrapped tea.Msg
}

type msgToSelectedNodeView struct {
	wrapped tea.Msg
}

type msgSelectableNodesUpdated struct {
	selectableNodes []*graph.Node
}
