package mainscreen

import (
	"pkb-agent/graph"

	tea "github.com/charmbracelet/bubbletea"
)

type MsgGraphLoaded struct {
	graph *graph.Graph
}

type MsgUpdateNodeList struct{}

type MsgToSelectableNodeList struct {
	wrapped tea.Msg
}

type MsgToSelectedNodeList struct {
	wrapped tea.Msg
}
