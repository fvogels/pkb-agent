package mainscreen

import (
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/components/textinput"
	"pkb-agent/ui/debug"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"slices"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	graph *graph.Graph
	size  util.Size
	mode  mode

	remainingNodes []*graph.Node
	selectedNodes  []*graph.Node

	remainingNodeView listview.Model[*graph.Node]
	selectedNodeView  listview.Model[*graph.Node]
	textInput         textinput.Model
}

func New() Model {
	renderer := func(node *graph.Node) string {
		return node.Name
	}

	return Model{
		mode:              viewMode{},
		remainingNodeView: listview.New(renderer, true),
		selectedNodeView:  listview.New(renderer, true),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Sequence(
		model.remainingNodeView.Init(),
		model.signalLoadGraph(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.KeyMsg:
		return model.onKeyPressed(message)

	case tea.WindowSizeMsg:
		return model.onResized(message)

	case MsgGraphLoaded:
		return model.onGraphLoaded(message)

	case textinput.MsgInputUpdated:
		return model.onInputUpdated(message)

	case msgToRemainingNodeView:
		updatedRemainingNodeList, command := model.remainingNodeView.TypedUpdate(message.wrapped)
		model.remainingNodeView = updatedRemainingNodeList
		return model, command

	case msgToSelectedNodeView:
		updatedSelectedNodeList, command := model.selectedNodeView.TypedUpdate(message.wrapped)
		model.selectedNodeView = updatedSelectedNodeList
		return model, command

	case msgRemainingNodesUpdated:
		model.remainingNodes = message.remainingNodes
		updatedRemainingNodesView, command := model.remainingNodeView.TypedUpdate(
			listview.MsgSetItems[*graph.Node]{
				Items: &SliceAdapter[*graph.Node]{
					slice: model.remainingNodes,
				},
			},
		)
		model.remainingNodeView = updatedRemainingNodesView
		return model, command

	default:
		updatedRemainingNodeView, command1 := model.remainingNodeView.TypedUpdate(message)
		model.remainingNodeView = updatedRemainingNodeView

		updatedSelectedNodeView, command2 := model.selectedNodeView.TypedUpdate(message)
		model.selectedNodeView = updatedSelectedNodeView

		updatedTextInput, command3 := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, tea.Batch(command1, command2, command3)
	}
}

func (model Model) onInputUpdated(_ textinput.MsgInputUpdated) (Model, tea.Cmd) {
	return model, model.signalUpdateRemainingNodes()
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	return model.mode.onKeyPressed(model, message)
}

func (model Model) View() string {
	return lipgloss.JoinVertical(
		0,
		lipgloss.NewStyle().Height(5).Render(model.selectedNodeView.View()),
		lipgloss.NewStyle().Height(model.size.Height-6).Render(model.remainingNodeView.View()),
		model.textInput.View(),
	)
}

func (model Model) onGraphLoaded(message MsgGraphLoaded) (Model, tea.Cmd) {
	model.graph = message.graph
	return model, model.signalUpdateRemainingNodes()
}

func (model *Model) signalLoadGraph() tea.Cmd {
	return func() tea.Msg {
		return MsgGraphLoaded{
			graph: loadGraph(),
		}
	}
}

func loadGraph() *graph.Graph {
	slog.Debug("loading graph")

	loader := metaloader.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	g, err := graph.LoadGraph(path, loader)
	if err != nil {
		panic("failed to load graph!")
	}

	return g
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	updatedSelectedNodeView, command1 := model.selectedNodeView.TypedUpdate(tea.WindowSizeMsg{
		Width:  message.Width,
		Height: 5,
	})
	model.selectedNodeView = updatedSelectedNodeView

	updatedRemainingNodeView, command2 := model.remainingNodeView.TypedUpdate(tea.WindowSizeMsg{
		Width:  message.Width,
		Height: message.Height - 6,
	})
	model.remainingNodeView = updatedRemainingNodeView

	return model, tea.Batch(command1, command2)
}

func (model Model) signalUpdateRemainingNodes() tea.Cmd {
	input := model.textInput.GetInput()
	iterator := model.graph.FindMatchingNodes(input)
	selectedNodes := model.selectedNodes

	return func() tea.Msg {
		nameSet := util.NewSet[string]()
		remaining := []*graph.Node{}

		for iterator.Current() != nil {
			// The same node can occur more than once during iteration
			// Ensure that we add each node only once to remainingNodes
			name := iterator.Current().Name
			if nameSet.Contains(name) {
				iterator.Next()
				continue
			}

			if !util.All(selectedNodes, func(selectedNode *graph.Node) bool {
				return slices.Contains(iterator.Current().Links, selectedNode.Name)
			}) {
				iterator.Next()
				continue
			}

			nameSet.Add(name)
			remaining = append(remaining, iterator.Current())
			iterator.Next()
		}

		sort.Slice(remaining, func(i, j int) bool {
			return remaining[i].Name < remaining[j].Name
		})
		return msgRemainingNodesUpdated{
			remainingNodes: remaining,
		}
	}
}

type SliceAdapter[T any] struct {
	slice []T
}

func NewSliceAdapter[T any](slice []T) *SliceAdapter[T] {
	return &SliceAdapter[T]{
		slice: slice,
	}
}

func (adapter *SliceAdapter[T]) Length() int {
	return len(adapter.slice)
}

func (adapter *SliceAdapter[T]) At(index int) T {
	return adapter.slice[index]
}
