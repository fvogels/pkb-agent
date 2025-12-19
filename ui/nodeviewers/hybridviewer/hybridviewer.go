package hybridviewer

import (
	"log/slog"
	"pkb-agent/graph/nodes/hybrid"
	"pkb-agent/ui/components/markdownview"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size      util.Size
	nodeExtra *hybrid.Extra
	nodeData  *hybrid.Data
	viewer    markdownview.Model
}

func New(nodeData *hybrid.Extra) Model {
	return Model{
		nodeExtra: nodeData,
		viewer:    markdownview.New(),
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadNodeData()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) UpdateViewer(message tea.Msg) (nodeviewers.Viewer, tea.Cmd) {
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

func (model *Model) signalLoadNodeData() tea.Cmd {
	extra := model.nodeExtra

	return func() tea.Msg {
		data, err := extra.GetData()
		if err != nil {
			slog.Debug("Error whlie reading node data", slog.String("error", err.Error()))
			panic("failed to load node's data")
		}

		return msgMarkdownLoaded{
			data,
		}
	}
}

func (model Model) onMarkdownLoaded(message msgMarkdownLoaded) (Model, tea.Cmd) {
	model.nodeData = message.data

	commands := []tea.Cmd{}
	if len(model.nodeData.MarkdownSource) > 0 {
		util.UpdateChild(&model.viewer, markdownview.MsgSetSource{
			Source: model.nodeData.MarkdownSource,
		}, &commands)
	}

	return model, tea.Batch(commands...)
}

func (model Model) GetKeyBindings() []key.Binding {
	return []key.Binding{}
}
