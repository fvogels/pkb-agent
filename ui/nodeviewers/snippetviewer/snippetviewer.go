package snippetviewer

import (
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/ui/debug"
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
		source:   "loading",
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadSnippet()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case msgSnippetLoaded:
		return model.onSnippetLoaded(message)
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

func (model *Model) signalLoadSnippet() tea.Cmd {
	data := model.nodeData

	return func() tea.Msg {
		source, err := data.GetSource()
		if err != nil {
			panic("failed to get snippet source")
		}

		return msgSnippetLoaded{
			source: source,
		}
	}
}

func (model Model) onSnippetLoaded(message msgSnippetLoaded) (Model, tea.Cmd) {
	model.source = message.source
	return model, nil
}
