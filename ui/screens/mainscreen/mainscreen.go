package mainscreen

import (
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	"pkb-agent/ui/components/nodeselectionview"
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

	nodeSelectionView nodeselectionview.Model
	textInput         textinput.Model
}

func New() Model {
	return Model{
		mode:              viewMode{},
		nodeSelectionView: nodeselectionview.New(),
		textInput:         textinput.New(),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.nodeSelectionView.Init(),
		model.textInput.Init(),
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

	case msgRemainingNodesUpdated:
		model.remainingNodes = message.remainingNodes

		return util.UpdateSingleChild(&model, &model.nodeSelectionView, nodeselectionview.MsgSetRemainingNodes{
			RemainingNodes: &SliceAdapter[*graph.Node]{
				slice: model.remainingNodes,
			},
		})

	default:
		commands := []tea.Cmd{}

		util.UpdateChild(&model.nodeSelectionView, message, &commands)
		util.UpdateChild(&model.textInput, message, &commands)

		return model, tea.Batch(commands...)
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
		lipgloss.NewStyle().Height(model.size.Height-1).Render(model.nodeSelectionView.View()),
		model.mode.renderStatusBar(&model),
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
		slog.Debug("error loading graph", slog.String("error", err.Error()))
		panic("failed to load graph!")
	}

	return g
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	commands := []tea.Cmd{}
	util.UpdateChild(&model.nodeSelectionView, tea.WindowSizeMsg{
		Width:  message.Width,
		Height: message.Height - 1,
	}, &commands)
	util.UpdateChild(&model.textInput, tea.WindowSizeMsg{
		Width:  message.Width,
		Height: 1,
	}, &commands)

	return model, tea.Batch(commands...)
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

func (model Model) onSelectPreviousRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(&model, &model.nodeSelectionView, nodeselectionview.MsgSelectPrevious{})
}

func (model Model) onSelectNextRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(&model, &model.nodeSelectionView, nodeselectionview.MsgSelectNext{})
}

func (model Model) onSelectNode() (Model, tea.Cmd) {
	model.mode = viewMode{}

	selectedNode := model.nodeSelectionView.GetSelectedRemainingNode()
	if selectedNode != nil {
		updatedSelectedNodes := append(model.selectedNodes, selectedNode)

		commands := []tea.Cmd{}
		util.UpdateChild(&model.textInput, textinput.MsgClear{}, &commands)
		updatedModel, command := model.setSelectedNodes(updatedSelectedNodes)
		model = updatedModel
		commands = append(commands, command)

		util.UpdateChild(&model.nodeSelectionView, nodeselectionview.MsgSetSelectedNodes{
			SelectedNodes: NewSliceAdapter(updatedSelectedNodes),
		}, &commands)

		commands = append(commands, model.signalUpdateRemainingNodes())

		return model, tea.Batch(
			commands...,
		)
	}

	return model, nil
}

func (model Model) onUnselectLast() (Model, tea.Cmd) {
	if len(model.selectedNodes) > 0 {
		updatedSelectedNodes := model.selectedNodes[:len(model.selectedNodes)-1]

		return model.setSelectedNodes(updatedSelectedNodes)
	} else {
		return model, nil
	}
}

func (model Model) setSelectedNodes(selectedNodes []*graph.Node) (Model, tea.Cmd) {
	model.selectedNodes = selectedNodes

	commands := []tea.Cmd{}
	util.UpdateChild(&model.nodeSelectionView, nodeselectionview.MsgSetSelectedNodes{
		SelectedNodes: NewSliceAdapter(model.selectedNodes),
	}, &commands)

	commands = append(commands, model.signalUpdateRemainingNodes())

	return model, tea.Batch(commands...)
}
