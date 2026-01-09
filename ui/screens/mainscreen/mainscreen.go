package mainscreen

import (
	"log/slog"
	"pkb-agent/graph/loaders/sequence"
	"pkb-agent/graph/node"
	"pkb-agent/pkg"
	"pkb-agent/ui/components/helpbar"
	"pkb-agent/ui/components/nodeselectionview"
	"pkb-agent/ui/components/nodeviewer"
	"pkb-agent/ui/components/textinput"
	"pkb-agent/ui/debug"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"slices"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	graph                    *pkg.Graph
	size                     util.Size
	mode                     mode
	includeLinkedNodes       bool
	includeIndirectAncestors bool
	nodeViewerKeyBindings    []key.Binding

	remainingNodes []*pkg.Node
	selectedNodes  []*pkg.Node

	nodeSelectionView nodeselectionview.Model
	nodeViewer        nodeviewer.Model
	textInput         textinput.Model
	helpBar           helpbar.Model

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

	createUpdateKeyBindingsMessage := func(keyBindings []key.Binding) tea.Msg {
		return node.MsgUpdateNodeViewerBindings{
			KeyBindings: keyBindings,
		}
	}

	model := Model{
		mode:                     viewMode,
		includeLinkedNodes:       true,
		includeIndirectAncestors: true,
		nodeSelectionView:        nodeselectionview.New(),
		textInput:                textinput.New(),
		helpBar:                  helpbar.New(),
		nodeViewer:               nodeviewer.New(createUpdateKeyBindingsMessage),
		layoutConfiguration:      &layoutConfiguration,
		viewMode:                 viewMode,
		inputMode:                inputMode,
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.nodeSelectionView.Init(),
		model.textInput.Init(),
		model.signalLoadGraph(),
		model.signalActivateMode(),
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

	case msgGraphLoaded:
		return model.onGraphLoaded(message)

	case textinput.MsgInputUpdated:
		return model.onInputUpdated(message)

	case msgRemainingNodesDetermined:
		model.remainingNodes = message.remainingNodes

		commands := []tea.Cmd{}

		util.UpdateChild(&model.nodeSelectionView, nodeselectionview.MsgSetRemainingNodes{
			RemainingNodes: &SliceAdapter[*pkg.Node]{
				slice: model.remainingNodes,
			},
			SelectionIndex: message.selectionIndex,
		}, &commands)

		util.UpdateChild(&model.helpBar, helpbar.MsgSetNodeCounts{
			Total:     model.graph.GetNodeCount(),
			Remaining: len(model.remainingNodes),
		}, &commands)

		return model, tea.Batch(commands...)

	case nodeselectionview.MsgRemainingNodeHighlighted:
		return model.onNodeHighlighted(message)

	case msgActivateMode:
		command := model.mode.activate(&model)
		return model, command

	case node.MsgUpdateNodeViewerBindings:
		model.nodeViewerKeyBindings = message.KeyBindings
		return model.refreshHelpBar()

	case msgSwitchMode:
		model.mode = message.mode
		command1 := model.mode.activate(&model)
		var command2 tea.Cmd
		model, command2 = model.refreshHelpBar()
		return model, tea.Batch(command1, command2)

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

func (model Model) onGraphLoaded(message msgGraphLoaded) (Model, tea.Cmd) {
	model.graph = message.graph
	return model, model.signalRefreshRemainingNodes(false)
}

func (model *Model) signalLoadGraph() tea.Cmd {
	return func() tea.Msg {
		g, err := loadGraph()
		if err != nil {
			strings.Lines(err.Error())(func(errorMessage string) bool {
				slog.Error("Failed to load graph", slog.String("error", strings.TrimSpace(errorMessage)))
				return true
			})

			panic("Failed to load graph")
		}

		return msgGraphLoaded{
			graph: g,
		}
	}
}

func loadGraph() (*pkg.Graph, error) {
	loader := sequence.New()
	path := pathlib.New(`F:\repos\pkb\pkb-data\root.yaml`)

	g, err := pkg.LoadGraph(path, loader)
	if err != nil {
		return nil, err
	}

	return g, nil
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
	highlightedNode := model.nodeSelectionView.GetSelectedRemainingNode()

	if len(input) == 0 {
		keepSameNodeHighlighted = true
	}

	return func() tea.Msg {
		remainingNodes := determineRemainingNodes(
			input,
			model.graph,
			selectedNodes,
			model.includeLinkedNodes,
			model.includeIndirectAncestors,
		)

		// Probably redundant; look into it
		sort.Slice(remainingNodes, func(i, j int) bool {
			return strings.ToLower(remainingNodes[i].GetName()) < strings.ToLower(remainingNodes[j].GetName())
		})

		highlightIndex := 0
		var target string
		if !keepSameNodeHighlighted || highlightedNode == nil {
			target = input
		} else {
			target = strings.ToLower(highlightedNode.GetName())
		}

		bestMatchIndex, found := slices.BinarySearchFunc(
			remainingNodes,
			target,
			func(node *pkg.Node, target string) int {
				nodeName := strings.ToLower(node.GetName())
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

func (model Model) setSelectedNodes(selectedNodes []*pkg.Node) (Model, tea.Cmd) {
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
	highlightedNode := message.Node

	// Reset key bindings, the specialized node viewer will update it again later
	model.nodeViewerKeyBindings = nil

	var shownNode *pkg.Node

	if highlightedNode == nil {
		// No node was highlighted, so take the last selected node

		if len(model.selectedNodes) == 0 {
			// Should not occur, but handle gracefully
			return model, nil
		}

		shownNode = model.selectedNodes[len(model.selectedNodes)-1]
	} else {
		shownNode = highlightedNode
	}

	var command1 tea.Cmd
	var command2 tea.Cmd
	model, command1 = model.showNode(shownNode)
	model, command2 = model.refreshHelpBar()

	return model, tea.Batch(command1, command2)
}

func (model Model) showNode(node *pkg.Node) (Model, tea.Cmd) {
	commands := []tea.Cmd{}

	util.UpdateChild(&model.nodeViewer, nodeviewer.MsgSetNode{Node: node}, &commands)

	return model, tea.Batch(commands...)
}

func (model Model) refreshHelpBar() (Model, tea.Cmd) {
	commands := []tea.Cmd{}

	util.UpdateChild(&model.helpBar, helpbar.MsgSetKeyBindings{
		KeyBindings: append(model.mode.getKeyBindings(), model.nodeViewerKeyBindings...),
	}, &commands)

	return model, tea.Batch(commands...)
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

func (model Model) toggleIncludeIndirectAncestors() (Model, tea.Cmd) {
	model.includeIndirectAncestors = !model.includeIndirectAncestors

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

func (model Model) signalActivateMode() tea.Cmd {
	return func() tea.Msg {
		return msgActivateMode{}
	}
}
