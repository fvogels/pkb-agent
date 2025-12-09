package mainscreen

import (
	tea "github.com/charmbracelet/bubbletea"
)

type viewMode struct{}

func (mode viewMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "q":
		return model, tea.Quit

	case "s":
		model.mode = inputMode{}
		return model, nil

	case "down":
		return model.onSelecNextRemainingNode()

	case "up":
		return model.onSelectPreviousRemainingNode()

	case "enter":
		return model.onSelectNode()

	default:
		return model, nil
	}
}

func (mode viewMode) renderStatusBar(model *Model) string {
	return ""
}
