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
	graph *graph.Graph
	size  util.Size
	mode  mode

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
		nodeSelectionViewHeight: 10,
	}
	viewMode := NewViewMode(&layoutConfiguration)
	inputMode := NewInputMode(&layoutConfiguration)

	model := Model{
		mode:                viewMode,
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

	case msgRemainingNodesUpdated:
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
	return model, model.signalUpdateRemainingNodes()
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	return model.mode.onKeyPressed(model, message)
}

func (model Model) View() string {
	return model.mode.render(&model)
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

	command := model.mode.resize(&model, model.size)

	return model, command
}

func (model Model) signalUpdateRemainingNodes() tea.Cmd {
	input := strings.ToLower(model.textInput.GetInput())
	iterator := model.graph.FindMatchingNodes(input)
	selectedNodes := model.selectedNodes

	return func() tea.Msg {
		// nameSet is used to prevent duplicates
		// Adding the selected nodes ensures that already selected nodes do not appear as remaining choices
		nameSet := util.NewSetFromSlice(util.Map(selectedNodes, func(node *graph.Node) string { return node.Name }))
		remaining := []*graph.Node{}

		// We need to keep track of it so that we can have it selected in the list
		var bestMatch *string = nil

		for iterator.Current() != nil {
			// The same node can occur more than once during iteration
			// Ensure that we add each node only once to remainingNodes
			name := iterator.Current().Name
			if nameSet.Contains(name) {
				iterator.Next()
				continue
			}

			if bestMatch == nil {
				bestMatch = &name
			} else {
				slog.Debug("!!!")
				if strings.HasPrefix(strings.ToLower(name), input) {
					if len(name) < len(*bestMatch) {
						*bestMatch = name
					}
				}
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

		imax := len(remaining)
		for i := 0; i != imax; i++ {
			node := remaining[i]

			for _, linkedNodeName := range node.Links {
				if !nameSet.Contains(linkedNodeName) {
					nameSet.Add(linkedNodeName)
					linkedNode := model.graph.FindNode(linkedNodeName)
					remaining = append(remaining, linkedNode)
				}
			}
		}

		sort.Slice(remaining, func(i, j int) bool {
			return strings.ToLower(remaining[i].Name) < strings.ToLower(remaining[j].Name)
		})

		if bestMatch != nil {
			slog.Debug("best match", "value", *bestMatch)
		}

		bestMatchIndex := 0
		if bestMatch != nil {
			var found bool
			bestMatchIndex, found = slices.BinarySearchFunc(
				remaining,
				*bestMatch,
				func(node *graph.Node, target string) int {
					if node.Name < target {
						return -1
					}
					if node.Name > target {
						return 1
					}
					return 0
				},
			)

			if !found {
				bestMatchIndex = 0
			}
		}

		slog.Debug("updating remaining nodes", slog.Int("bestMatchIndex", bestMatchIndex))
		return msgRemainingNodesUpdated{
			remainingNodes: remaining,
			selectionIndex: bestMatchIndex,
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
