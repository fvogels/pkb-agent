package bookmarkviewer

import (
	"pkb-agent/extern"
	"pkb-agent/graph/nodes/bookmark"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var keyMap = struct {
	OpenInBrowser key.Binding
}{
	OpenInBrowser: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open"),
	),
}

type Model struct {
	size                           util.Size
	nodeData                       *bookmark.Extra
	createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg
}

func New(createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg, nodeData *bookmark.Extra) Model {
	return Model{
		nodeData:                       nodeData,
		createUpdateKeyBindingsMessage: createUpdateKeyBindingsMessage,
	}
}

func (model Model) Init() tea.Cmd {
	return func() tea.Msg {
		return model.createUpdateKeyBindingsMessage([]key.Binding{
			keyMap.OpenInBrowser,
		})
	}
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

	case tea.KeyMsg:
		return model.onKeyPressed(message)
	}

	return model, nil
}

func (model Model) View() string {
	return model.nodeData.Description
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, keyMap.OpenInBrowser):
		return model.onOpenInBrowser()

	default:
		return model, nil
	}
}

func (model Model) onOpenInBrowser() (Model, tea.Cmd) {
	if err := extern.OpenURLInBrowser(model.nodeData.URL); err != nil {
		panic("failed to open browser")
	}

	return model, nil
}
