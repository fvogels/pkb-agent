package mainscreen

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var viewModeKeyMap = struct {
	Quit              key.Binding
	SwitchToInputMode key.Binding
	Next              key.Binding
	Previous          key.Binding
	Select            key.Binding
}{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	SwitchToInputMode: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "search"),
	),
	Next: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next"),
	),
	Previous: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "next"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
}

type viewMode struct{}

func (mode viewMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, viewModeKeyMap.Quit):
		return model, tea.Quit

	case key.Matches(message, viewModeKeyMap.SwitchToInputMode):
		model.mode = inputMode{}
		return model, nil

	case key.Matches(message, viewModeKeyMap.Next):
		return model.onSelecNextRemainingNode()

	case key.Matches(message, viewModeKeyMap.Previous):
		return model.onSelectPreviousRemainingNode()

	case key.Matches(message, viewModeKeyMap.Select):
		return model.onSelectNode()

	default:
		return model, nil
	}
}

func (mode viewMode) renderStatusBar(model *Model) string {
	return ""
}
