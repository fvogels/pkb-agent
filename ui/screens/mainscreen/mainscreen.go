package mainscreen

import (
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	"pkb-agent/ui/components/nodeselectionview"
	"pkb-agent/ui/components/textinput"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers/nodeviewer"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"slices"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	graph              *graph.Graph
	size               util.Size
	mode               mode
	includeLinkedNodes bool

	remainingNodes []*graph.Node
	selectedNodes  []*graph.Node

	nodeSelectionView nodeselectionview.Model
	nodeViewer        nodeviewer.Model
	textInput         textinput.Model

	layoutConfiguration *layoutConfiguration

	viewMode  *viewMode
	inputMode *inputMode
}

func New() Model {
	layoutConfiguration := layoutConfiguration{
		nodeSelectionViewHeight: 20,
	}
	viewMode := NewViewMode(&layoutConfiguration)
	inputMode := NewInputMode(&layoutConfiguration)

	model := Model{
		mode:                viewMode,
		includeLinkedNodes:  true,
		nodeSelectionView:   nodeselectionview.New(),
		textInput:           textinput.New(),
		nodeViewer:          nodeviewer.New(),
		layoutConfiguration: &layoutConfiguration,
		viewMode:            viewMode,
		inputMode:           inputMode,
	}

	return model
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

	case msgRemainingNodesDetermined:
		model.remainingNodes = message.remainingNodes

		return util.UpdateSingleChild(&model, &model.nodeSelectionView, nodeselectionview.MsgSetRemainingNodes{
			RemainingNodes: &SliceAdapter[*graph.Node]{
				slice: model.remainingNodes,
			},
			SelectionIndex: message.selectionIndex,
		})

	case nodeselectionview.MsgRemainingNodeHighlighted:
		return model.onNodeHighlighted(message)

	default:
		commands := []tea.Cmd{}

		util.UpdateChild(&model.nodeSelectionView, message, &commands)
		util.UpdateChild(&model.textInput, message, &commands)
		util.UpdateChild(&model.nodeViewer, message, &commands)

		return model, tea.Batch(commands...)
	}
}

func (model Model) onInputUpdated(_ textinput.MsgInputUpdated) (Model, tea.Cmd) {
	return model, model.signalRefreshRemainingNodes(false)
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	return model.mode.onKeyPressed(model, message)
}

func (model Model) View() string {
	return model.mode.render(&model)
}

func (model Model) onGraphLoaded(message MsgGraphLoaded) (Model, tea.Cmd) {
	model.graph = message.graph
	return model, model.signalRefreshRemainingNodes(false)
}

func (model *Model) signalLoadGraph() tea.Cmd {
	return func() tea.Msg {
		return MsgGraphLoaded{
			graph: loadGraph(),
		}
	}
}

func loadGraph() *graph.Graph {
	loader := metaloader.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	g, err := graph.LoadGraph(path, loader)
	if err != nil {
		slog.Error("error loading graph", slog.String("error", err.Error()))
		panic("failed to load graph!")
	}

	return g
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	command := model.mode.resize(&model, model.size)

	return model, command
}

func (model Model) signalRefreshRemainingNodes(keepSameNodeHighlighted bool) tea.Cmd {
	input := strings.ToLower(model.textInput.GetInput())
	selectedNodes := model.selectedNodes
	highlighedNode := model.nodeSelectionView.GetSelectedRemainingNode()

	if len(input) == 0 {
		keepSameNodeHighlighted = true
	}

	return func() tea.Msg {
		remainingNodes := determineRemainingNodes(
			input,
			model.graph,
			selectedNodes,
			model.includeLinkedNodes,
		)

		sort.Slice(remainingNodes, func(i, j int) bool {
			return strings.ToLower(remainingNodes[i].Name) < strings.ToLower(remainingNodes[j].Name)
		})

		highlightIndex := 0
		var target string
		if !keepSameNodeHighlighted || highlighedNode == nil {
			target = input
		} else {
			target = strings.ToLower(highlighedNode.Name)
		}

		bestMatchIndex, found := slices.BinarySearchFunc(
			remainingNodes,
			target,
			func(node *graph.Node, target string) int {
				nodeName := strings.ToLower(node.Name)
				if strings.HasPrefix(nodeName, target) {
					return 0
				}
				if nodeName < target {
					return -1
				}
				return 1
			},
		)

		if !found {
			bestMatchIndex = 0
		}

		highlightIndex = bestMatchIndex

		return msgRemainingNodesDetermined{
			remainingNodes: remainingNodes,
			selectionIndex: highlightIndex,
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

func (model Model) onHighlightPreviousRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(
		&model,
		&model.nodeSelectionView,
		nodeselectionview.MsgHighlightRemainingNode{
			Index: model.nodeSelectionView.GetSelectedRemainingNodeIndex() - 1,
		},
	)
}

func (model Model) onHighlightNextRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(
		&model,
		&model.nodeSelectionView,
		nodeselectionview.MsgHighlightRemainingNode{
			Index: model.nodeSelectionView.GetSelectedRemainingNodeIndex() + 1,
		},
	)
}

func (model Model) onHighlightRemainingNodePageDown() (Model, tea.Cmd) {
	slog.Debug("!!!!")
	return util.UpdateSingleChild(
		&model,
		&model.nodeSelectionView,
		nodeselectionview.MsgHighlightRemainingNode{
			Index: model.nodeSelectionView.GetSelectedRemainingNodeIndex() + model.nodeSelectionView.GetRemaingNodesPageSize(),
		},
	)
}

func (model Model) onHighlightRemainingNodePageUp() (Model, tea.Cmd) {
	return util.UpdateSingleChild(
		&model,
		&model.nodeSelectionView,
		nodeselectionview.MsgHighlightRemainingNode{
			Index: model.nodeSelectionView.GetSelectedRemainingNodeIndex() - model.nodeSelectionView.GetRemaingNodesPageSize(),
		},
	)
}

func (model Model) onSelectFirstRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(&model, &model.nodeSelectionView, nodeselectionview.MsgHighlightRemainingNode{Index: 0})
}

func (model Model) onSelectLastRemainingNode() (Model, tea.Cmd) {
	return util.UpdateSingleChild(
		&model,
		&model.nodeSelectionView,
		nodeselectionview.MsgHighlightRemainingNode{
			Index: len(model.remainingNodes) - 1,
		},
	)
}

func (model Model) onNodeSelected() (Model, tea.Cmd) {
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

		commands = append(commands, model.signalRefreshRemainingNodes(false))

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

	commands = append(commands, model.signalRefreshRemainingNodes(false))

	return model, tea.Batch(commands...)
}

// onNodeHighlighted is called whenever a new node is being highlighted in the list of remaining nodes.
func (model Model) onNodeHighlighted(message nodeselectionview.MsgRemainingNodeHighlighted) (Model, tea.Cmd) {
	highlighedNode := message.Node

	if highlighedNode == nil {
		// No node was highlighted
		return model, nil
	} else {
		return util.UpdateSingleChild(&model, &model.nodeViewer, nodeviewer.MsgSetNode{Node: highlighedNode})
	}
}

func (model Model) updateLayoutConfiguration(update func(*layoutConfiguration)) (Model, tea.Cmd) {
	update(model.layoutConfiguration)
	command := model.mode.resize(&model, model.size)
	return model, command
}

func (model Model) toggleIncludeLinkedNodes() (Model, tea.Cmd) {
	model.includeLinkedNodes = !model.includeLinkedNodes

	return model, model.signalRefreshRemainingNodes(true)
}

func (model Model) setInput(input string) (Model, tea.Cmd) {
	return util.UpdateSingleChild(
		&model,
		&model.textInput,
		textinput.MsgSetInput{
			Input: input,
		},
	)
}
