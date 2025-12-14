package nodeselectionview

import (
	"log/slog"
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

	remainingNodesView := listview.New(renderer, true, wrapRemainingNodesViewMessage)

	selectedNodesView := listview.New(renderer, false, wrapSelectedNodesViewMessage)
	selectedNodesView.SetNonselectedStyle(lipgloss.NewStyle().Background(lipgloss.Color("#AAFFAA")))

	return Model{
		remainingNodes:     &emptyList{},
		selectedNodes:      &emptyList{},
		remainingNodesView: remainingNodesView,
		selectedNodesView:  selectedNodesView,
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

		return util.UpdateSingleChild(&model, &model.remainingNodesView, listview.MsgSetItems[*graph.Node]{
			Items:          message.RemainingNodes,
			SelectionIndex: message.SelectionIndex,
		})

	case MsgSetSelectedNodes:
		model.selectedNodes = message.SelectedNodes

		commands := []tea.Cmd{}
		util.UpdateChild(&model.selectedNodesView, listview.MsgSetItems[*graph.Node]{
			Items: message.SelectedNodes,
		}, &commands)

		// Number of selected nodes has changed
		// This affects selected node list's size
		updatedModel, command := model.updateChildSizes()
		commands = append(commands, command)

		return updatedModel, tea.Batch(commands...)

	case MsgSelectPrevious:
		return util.UpdateSingleChild(&model, &model.remainingNodesView, listview.MsgSelectPrevious{})

	case MsgSelectNext:
		return util.UpdateSingleChild(&model, &model.remainingNodesView, listview.MsgSelectNext{})

	case msgRemainingNodesWrapper:
		switch message := message.wrapped.(type) {
		case listview.MsgItemSelected[*graph.Node]:
			return model, model.signalRemainingNodeHighlighted(message.Item)

		case listview.MsgNoItemSelected:
			return model, model.signalRemainingNodeHighlighted(nil)

		default:
			slog.Warn("swallowed message from remaining nodes list")
			return model, nil
		}

	case msgSelectedNodesWrapper:
		slog.Warn("swallowed message from selected nodes list")
		return model, nil

	default:
		commands := []tea.Cmd{}

		util.UpdateChild(&model.selectedNodesView, message, &commands)
		util.UpdateChild(&model.remainingNodesView, message, &commands)

		return model, tea.Batch(commands...)
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
	commands := []tea.Cmd{}

	util.UpdateChild(&model.selectedNodesView, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.selectedNodes.Length(),
	}, &commands)

	util.UpdateChild(&model.remainingNodesView, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.size.Height - model.selectedNodes.Length(),
	}, &commands)

	return model, tea.Batch(commands...)
}

func (model Model) signalRemainingNodeHighlighted(node *graph.Node) tea.Cmd {
	return func() tea.Msg {
		return MsgRemainingNodeHighlighted{
			Node: node,
		}
	}
}

func wrapSelectedNodesViewMessage(message tea.Msg) tea.Msg {
	return msgSelectedNodesWrapper{
		wrapped: message,
	}
}

func wrapRemainingNodesViewMessage(message tea.Msg) tea.Msg {
	return msgRemainingNodesWrapper{
		wrapped: message,
	}
}
