package mainscreen

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type inputMode struct{}

var inputModeKeyMap = struct {
	Cancel   key.Binding
	Next     key.Binding
	Previous key.Binding
	Select   key.Binding
}{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
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

func (mode inputMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, inputModeKeyMap.Cancel):
		model.mode = viewMode{}
		return model, nil

	case key.Matches(message, viewModeKeyMap.Next):
		return model.onSelecNextRemainingNode()

	case key.Matches(message, viewModeKeyMap.Previous):
		return model.onSelectPreviousRemainingNode()

	case key.Matches(message, viewModeKeyMap.Select):
		return model.onSelectNode()

	default:
		updatedTextInput, command := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, command
	}
}

func (mode inputMode) renderStatusBar(model *Model) string {
	return model.textInput.View()
}
