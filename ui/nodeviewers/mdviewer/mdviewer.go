package mdviewer

import (
	"pkb-agent/graph/nodes/markdown"
	"pkb-agent/ui/components/markdownview"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size     util.Size
	nodeData *markdown.Extra
	viewer   markdownview.Model
}

func New(nodeData *markdown.Extra) Model {
	return Model{
		nodeData: nodeData,
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadMarkdown()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case msgMarkdownLoaded:
		return model.onMarkdownLoaded(message)

	default:
		return util.UpdateSingleChild(&model, &model.viewer, message)
	}
}

func (model Model) View() string {
	return model.viewer.View()
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return util.UpdateSingleChild(&model, &model.viewer, message)
}

func (model *Model) signalLoadMarkdown() tea.Cmd {
	data := model.nodeData

	return func() tea.Msg {
		source, err := data.GetSource()
		if err != nil {
			panic("failed to get markdown source")
		}

		return msgMarkdownLoaded{
			source,
		}
	}
}

func (model Model) onMarkdownLoaded(message msgMarkdownLoaded) (Model, tea.Cmd) {
	return util.UpdateSingleChild(&model, &model.viewer, markdownview.MsgSetSource{
		Source: message.source,
	})
}
