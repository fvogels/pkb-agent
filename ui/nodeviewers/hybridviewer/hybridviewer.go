package hybridviewer

import (
	"log/slog"
	"pkb-agent/extern"
	"pkb-agent/graph/nodes/hybrid"
	"pkb-agent/ui/components/markdownview"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size                           util.Size
	nodeExtra                      *hybrid.Extra
	nodeData                       *hybrid.Data
	viewer                         markdownview.Model
	createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg
}

var keyMap = struct {
	OpenLink key.Binding
}{
	OpenLink: key.NewBinding(
		key.WithKeys("w"),
		key.WithHelp("w", "www"),
	),
}

func New(createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg, nodeData *hybrid.Extra) Model {
	return Model{
		nodeExtra:                      nodeData,
		viewer:                         markdownview.New(),
		createUpdateKeyBindingsMessage: createUpdateKeyBindingsMessage,
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
		return model.onDataLoaded(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

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

func (model Model) onDataLoaded(message msgMarkdownLoaded) (Model, tea.Cmd) {
	model.nodeData = message.data

	commands := []tea.Cmd{model.signalUpdatedKeyBindings()}

	if len(model.nodeData.MarkdownSource) > 0 {
		util.UpdateChild(&model.viewer, markdownview.MsgSetSource{
			Source: model.nodeData.MarkdownSource,
		}, &commands)
	}

	keyMap.OpenLink.SetEnabled(len(model.nodeData.URL) > 0)

	return model, tea.Batch(commands...)
}

func (model Model) signalUpdatedKeyBindings() tea.Cmd {
	return func() tea.Msg {
		return model.createUpdateKeyBindingsMessage(model.determineKeyBindings())
	}
}

func (model Model) determineKeyBindings() []key.Binding {
	bindings := []key.Binding{}

	if keyMap.OpenLink.Enabled() {
		bindings = append(bindings, keyMap.OpenLink)
	}

	return bindings
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, keyMap.OpenLink):
		return model.onOpenURL()

	default:
		return model, nil
	}
}

func (model Model) onOpenURL() (Model, tea.Cmd) {
	if err := extern.OpenURLInBrowser(model.nodeData.URL); err != nil {
		panic("failed to open browser")
	}

	return model, nil
}
