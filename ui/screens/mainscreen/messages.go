package mainscreen

import (
	"pkb-agent/graph"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgGraphLoaded struct {
	graph *graph.Graph
}

type MsgUpdateNodeList struct{}

type msgToRemainingNodeView struct {
	wrapped tea.Msg
}

type msgToSelectedNodeView struct {
	wrapped tea.Msg
}

type msgRemainingNodesDetermined struct {
	remainingNodes []*graph.Node
	selectionIndex int
}

type msgActivateMode struct{}
