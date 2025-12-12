package bookmarkviewer

import (
	"pkb-agent/graph/nodes/bookmark"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size     util.Size
	nodeData *bookmark.Extra
}

func New(nodeData *bookmark.Extra) Model {
	return Model{
		nodeData: nodeData,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)
	}

	return model, nil
}

func (model Model) View() string {
	return model.nodeData.URL
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}
