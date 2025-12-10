package snippetviewer

import (
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size     util.Size
	nodeData *snippet.Extra
	source   string
}

func New(nodeData *snippet.Extra) Model {
	return Model{
		nodeData: nodeData,
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadSnippet
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case msgSnippetLoaded:
		model.source = message.source
	}

	return model, nil
}

func (model Model) View() string {
	return model.source
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}

func (model *Model) signalLoadSnippet() tea.Msg {
	return func() tea.Msg {
		return msgSnippetLoaded{
			source: "source!",
		}
	}
}
