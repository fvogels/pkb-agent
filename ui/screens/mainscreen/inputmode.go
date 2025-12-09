package mainscreen

import (
	"pkb-agent/graph"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/components/textinput"

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
		model.mode = viewMode{}

		if len(model.remainingNodes) > 0 {
			selectedNode := model.remainingNodeView.GetSelectedItem()
			model.selectedNodes = append(model.selectedNodes, selectedNode)

			updatedSelectedNodeView, command1 := model.selectedNodeView.TypedUpdate(listview.MsgSetItems[*graph.Node]{
				Items: NewSliceAdapter(model.selectedNodes),
			})
			model.selectedNodeView = updatedSelectedNodeView

			updatedTextInput, command2 := model.textInput.TypedUpdate(textinput.MsgClear{})
			model.textInput = updatedTextInput

			return model, tea.Batch(
				command1,
				command2,
				model.signalUpdateRemainingNodes(),
			)
		}

		return model, nil

	default:
		updatedTextInput, command := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, command
	}
}

func (mode inputMode) renderStatusBar(model *Model) string {
	return model.textInput.View()
}
