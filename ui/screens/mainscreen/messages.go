package mainscreen

import (
	"pkb-agent/pkg"

	tea "github.com/charmbracelet/bubbletea"
)

type msgGraphLoaded struct {
	graph *pkg.Graph
}

type MsgUpdateNodeList struct{}

type msgToRemainingNodeView struct {
	wrapped tea.Msg
}

type msgToSelectedNodeView struct {
	wrapped tea.Msg
}

type msgRemainingNodesDetermined struct {
	remainingNodes []*pkg.Node
	selectionIndex int
}

type msgActivateMode struct{}

type msgSwitchMode struct {
	mode mode
}
