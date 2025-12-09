package mainscreen

import (
	"pkb-agent/graph"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/components/textinput"

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
		updatedNodeList, command := model.remainingNodeView.TypedUpdate(listview.MsgSelectNext{})
		model.remainingNodeView = updatedNodeList
		return model, command

	case "up":
		updatedNodeList, command := model.remainingNodeView.TypedUpdate(listview.MsgSelectPrevious{})
		model.remainingNodeView = updatedNodeList
		return model, command

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
		return model, nil
	}
}

func (mode viewMode) renderStatusBar(model *Model) string {
	return ""
}
