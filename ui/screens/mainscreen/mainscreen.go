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
	graph         *graph.Graph
	size          util.Size
	selectedNodes []*graph.Node
	nodeList      listview.Model[NodeWrapper]
	// selectedNodeList listview.Model
	textInput textinput.Model
}

type NodeWrapper struct {
	*graph.Node
}

func (wrapper NodeWrapper) String() string {
	return wrapper.Name
}

func New() Model {
	return Model{
		nodeList: listview.New[NodeWrapper](true),
		// selectedNodeList: listview.New(true),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Sequence(
		model.nodeList.Init(),
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

	default:
		updatedNodeList, command1 := model.nodeList.TypedUpdate(message)
		model.nodeList = updatedNodeList

		updatedTextInput, command2 := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, tea.Batch(command1, command2)
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
		updatedNodeList, command := model.nodeList.TypedUpdate(listview.MsgSelectNext{})
		model.nodeList = updatedNodeList
		return model, command

	case "up":
		updatedNodeList, command := model.nodeList.TypedUpdate(listview.MsgSelectPrevious{})
		model.nodeList = updatedNodeList
		return model, command

	default:
		updatedTextInput, command := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, command
	}
}

func (model Model) View() string {
	return lipgloss.JoinVertical(
		0,
		lipgloss.NewStyle().Height(model.size.Height-1).Render(model.nodeList.View()),
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

	updatedNodeList, command := model.nodeList.TypedUpdate(tea.WindowSizeMsg{
		Width:  message.Width,
		Height: message.Height - 1,
	})
	model.nodeList = updatedNodeList

	return model, command
}

func (model Model) signalUpdateNodeList() tea.Cmd {
	return func() tea.Msg {
		input := model.textInput.GetInput()
		iterator := model.graph.FindNameMatches(input)
		nameTable := make(map[string]any)
		nodes := []NodeWrapper{}

		for iterator.Current() != nil {
			name := iterator.Current().Name
			if _, alreadyAdded := nameTable[name]; !alreadyAdded {
				nameTable[name] = nil
				nodes = append(nodes, NodeWrapper{Node: iterator.Current()})
			}
			iterator.Next()
		}

		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Name < nodes[j].Name
		})
		return listview.MsgSetItems[NodeWrapper]{
			Items: &SliceAdapter[NodeWrapper]{
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
