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
	nodeInfo                       *hybrid.Info
	nodeData                       *hybrid.Data
	viewer                         markdownview.Model
	actions                        []action
	createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg
}

type action struct {
	keyBinding key.Binding
	perform    func()
}

func New(createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg, nodeData *hybrid.Info) Model {
	return Model{
		nodeInfo:                       nodeData,
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
	info := model.nodeInfo

	return func() tea.Msg {
		data, err := info.GetData()
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
	model.actions = model.createCommands()
	commands := []tea.Cmd{model.signalUpdatedKeyBindings()}

	if len(model.nodeData.MarkdownSource) > 0 {
		util.UpdateChild(&model.viewer, markdownview.MsgSetSource{
			Source: model.nodeData.MarkdownSource,
		}, &commands)
	}

	return model, tea.Batch(commands...)
}

func (model Model) signalUpdatedKeyBindings() tea.Cmd {
	return func() tea.Msg {
		return model.createUpdateKeyBindingsMessage(model.determineKeyBindings())
	}
}

func (model Model) createCommands() []action {
	actions := []action{}
	keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	keyBindingIndex := 0

	for _, externalLink := range model.nodeData.ExternalLinks {
		binding := key.NewBinding(
			key.WithKeys(keys[keyBindingIndex]),
			key.WithHelp(keys[keyBindingIndex], externalLink.Description),
		)
		action := action{
			keyBinding: binding,
			perform:    func() { openURL(externalLink.URL) },
		}
		actions = append(actions, action)
		keyBindingIndex++
	}

	return actions
}

func (model Model) determineKeyBindings() []key.Binding {
	return util.Map(model.actions, func(action action) key.Binding {
		return action.keyBinding
	})
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	for _, action := range model.actions {
		if key.Matches(message, action.keyBinding) {
			action.perform()
			return model, nil
		}
	}

	return model, nil
}

func openURL(url string) {
	if err := extern.OpenURLInBrowser(url); err != nil {
		panic("failed to open browser")
	}
}
