package mainscreen

import (
	tea "github.com/charmbracelet/bubbletea"
)

type inputMode struct{}

func (mode inputMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "esc":
		model.mode = viewMode{}
		return model, nil

	case "down":
		return model.onSelecNextRemainingNode()

	case "up":
		return model.onSelectPreviousRemainingNode()

	case "enter":
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
