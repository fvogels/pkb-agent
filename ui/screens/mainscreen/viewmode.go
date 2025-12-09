package mainscreen

import (
	"pkb-agent/ui/components/listview"

	tea "github.com/charmbracelet/bubbletea"
)

type viewMode struct {
}

func (mode viewMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "q":
		return model, tea.Quit

	case " ":
		model.mode = inputMode{}
		return model, nil

	case "down":
		updatedNodeList, command := model.remainingNodeView.TypedUpdate(listview.MsgSelectNext{})
		model.remainingNodeView = updatedNodeList
		return model, command

	case "up":
		updatedNodeList, command := model.remainingNodeView.TypedUpdate(listview.MsgSelectPrevious{})
		model.remainingNodeView = updatedNodeList
		return model, command

	default:
		return model, nil
	}
}
