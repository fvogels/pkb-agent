package bookmark

import (
	"pkb-agent/extern"
	"pkb-agent/graph/node"
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
	size util.Size
	url  string
}

func NewViewer(url string) Model {
	return Model{
		url: url,
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalKeybindingsUpdate()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

	default:
		return model, nil
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
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
	if err := extern.OpenURLInBrowser(model.url); err != nil {
		panic("failed to open browser")
	}

	return model, nil
}

func (model Model) View() string {
	return model.url
}

func (model Model) signalKeybindingsUpdate() tea.Cmd {
	return func() tea.Msg {
		return node.MsgUpdateNodeViewerBindings{
			KeyBindings: []key.Binding{
				keyMap.OpenInBrowser,
			},
		}
	}
}
