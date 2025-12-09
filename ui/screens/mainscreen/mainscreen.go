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
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	graph *graph.Graph
	size  util.Size

	selectableNodes []*graph.Node
	selectedNodes   []*graph.Node

	selectableNodeList listview.Model[*graph.Node]
	selectedNodeList   listview.Model[*graph.Node]
	textInput          textinput.Model
}

func New() Model {
	renderer := func(node *graph.Node) string {
		return node.Name
	}

	return Model{
		selectableNodeList: listview.New(renderer, true),
		selectedNodeList:   listview.New(renderer, true),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Sequence(
		model.selectableNodeList.Init(),
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

	case MsgToSelectableNodeList:
		updatedSelectableNodeList, command := model.selectableNodeList.TypedUpdate(message)
		model.selectableNodeList = updatedSelectableNodeList
		return model, command

	case MsgToSelectedNodeList:
		updatedSelectedNodeList, command := model.selectedNodeList.TypedUpdate(message)
		model.selectedNodeList = updatedSelectedNodeList
		return model, command

	default:
		updatedSelectableNodeList, command1 := model.selectableNodeList.TypedUpdate(message)
		model.selectableNodeList = updatedSelectableNodeList

		updatedSelectedNodeList, command2 := model.selectedNodeList.TypedUpdate(message)
		model.selectedNodeList = updatedSelectedNodeList

		updatedTextInput, command3 := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, tea.Batch(command1, command2, command3)
	}
}

func (model Model) onInputUpdated(_ textinput.MsgInputUpdated) (Model, tea.Cmd) {
	return model, model.signalUpdateNodeList()
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "esc":
		return model, tea.Quit

	case "down":
		updatedNodeList, command := model.selectableNodeList.TypedUpdate(listview.MsgSelectNext{})
		model.selectableNodeList = updatedNodeList
		return model, command

	case "up":
		updatedNodeList, command := model.selectableNodeList.TypedUpdate(listview.MsgSelectPrevious{})
		model.selectableNodeList = updatedNodeList
		return model, command

	case "enter":
		selectedNode := model.selectableNodeList.GetSelectedItem()
		model.selectedNodes = append(model.selectableNodes, selectedNode)
		return model, model.signalUpdateNodeList()

	default:
		updatedTextInput, command := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, command
	}
}

func (model Model) View() string {
	return lipgloss.JoinVertical(
		0,
		lipgloss.NewStyle().Height(5).Render(model.selectedNodeList.View()),
		lipgloss.NewStyle().Height(model.size.Height-6).Render(model.selectableNodeList.View()),
		model.textInput.View(),
	)
}

func (model Model) onGraphLoaded(message MsgGraphLoaded) (Model, tea.Cmd) {
	model.graph = message.graph
	return model, model.signalUpdateNodeList()
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

	updatedNodeList, command := model.selectableNodeList.TypedUpdate(tea.WindowSizeMsg{
		Width:  message.Width,
		Height: message.Height - 1,
	})
	model.selectableNodeList = updatedNodeList

	return model, command
}

func (model Model) signalUpdateNodeList() tea.Cmd {
	return func() tea.Msg {
		input := model.textInput.GetInput()
		iterator := model.graph.FindNameMatches(input)
		nameSet := util.NewSet[string]()
		nodes := []*graph.Node{}

		for iterator.Current() != nil {
			name := iterator.Current().Name
			if !nameSet.Contains(name) {
				nameSet.Add(name)
				nodes = append(nodes, iterator.Current())
			}
			iterator.Next()
		}

		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})
		return listview.MsgSetItems[*graph.Node]{
			Items: &SliceAdapter[*graph.Node]{
				slice: nodes,
			},
		}
	}
}

type SliceAdapter[T any] struct {
	slice []T
}

func (adapter *SliceAdapter[T]) Length() int {
	return len(adapter.slice)
}

func (adapter *SliceAdapter[T]) At(index int) T {
	return adapter.slice[index]
}
