package nodeselectionview

import (
	"pkb-agent/graph"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size util.Size

	remainingNodes List
	selectedNodes  List

	remainingNodesView listview.Model[*graph.Node]
	selectedNodesView  listview.Model[*graph.Node]
}

type List interface {
	At(index int) *graph.Node
	Length() int
}

type emptyList struct{}

func (emptyList *emptyList) At(index int) *graph.Node {
	panic("invalid operation")
}

func (emptyList *emptyList) Length() int {
	return 0
}

func New() Model {
	renderer := func(node *graph.Node) string {
		return node.Name
	}

	return Model{
		remainingNodes:     &emptyList{},
		selectedNodes:      &emptyList{},
		remainingNodesView: listview.New(renderer, true),
		selectedNodesView:  listview.New(renderer, false),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.remainingNodesView.Init(),
		model.selectedNodesView.Init(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetRemainingNodes:
		model.remainingNodes = message.RemainingNodes
		updatedRemainingNodesView, command := model.remainingNodesView.TypedUpdate(listview.MsgSetItems[*graph.Node]{
			Items: message.RemainingNodes,
		})
		model.remainingNodesView = updatedRemainingNodesView
		return model, command

	case MsgSetSelectedNodes:
		model.selectedNodes = message.SelectedNodes
		updatedSelectedNodesView, command1 := model.selectedNodesView.TypedUpdate(listview.MsgSetItems[*graph.Node]{
			Items: message.SelectedNodes,
		})
		model.selectedNodesView = updatedSelectedNodesView

		// Number of selected nodes has changed
		// This affects selected node list's size
		updatedModel, command2 := model.updateChildSizes()

		return updatedModel, tea.Batch(command1, command2)

	case MsgSelectPrevious:
		updatedRemainingNodesView, command := model.remainingNodesView.TypedUpdate(listview.MsgSelectPrevious{})
		model.remainingNodesView = updatedRemainingNodesView
		return model, command

	case MsgSelectNext:
		updatedRemainingNodesView, command := model.remainingNodesView.TypedUpdate(listview.MsgSelectNext{})
		model.remainingNodesView = updatedRemainingNodesView
		return model, command

	default:
		updatedSelectedNodesView, command1 := model.selectedNodesView.TypedUpdate(message)
		model.selectedNodesView = updatedSelectedNodesView

		updatedRemainingNodesView, command2 := model.remainingNodesView.TypedUpdate(message)
		model.remainingNodesView = updatedRemainingNodesView

		return model, tea.Batch(command1, command2)
	}
}

func (model Model) View() string {
	children := []string{}

	if model.selectedNodes.Length() > 0 {
		children = append(children, model.selectedNodesView.View())
	}

	children = append(children, model.remainingNodesView.View())

	return lipgloss.JoinVertical(0, children...)
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model.updateChildSizes()
}

func (model Model) GetSelectedRemainingNode() *graph.Node {
	return model.remainingNodesView.GetSelectedItem()
}

func (model Model) updateChildSizes() (Model, tea.Cmd) {
	updatedSelectedNodesView, command1 := model.selectedNodesView.TypedUpdate(tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.selectedNodes.Length(),
	})
	model.selectedNodesView = updatedSelectedNodesView

	updatedRemainingNodesView, command2 := model.remainingNodesView.TypedUpdate(tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.size.Height - model.selectedNodes.Length(),
	})
	model.remainingNodesView = updatedRemainingNodesView

	return model, tea.Batch(command1, command2)
}
