package mainscreen

import (
	"log/slog"
	"pkb-agent/graph"
	"pkb-agent/graph/metaloader"
	"pkb-agent/ui/components/listview"
	"pkb-agent/ui/debug"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	graph    *graph.Graph
	nodeList listview.Model
	size     util.Size
}

func New() Model {
	debug.Milestone()

	return Model{
		nodeList: listview.New(true),
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Sequence(
		model.nodeList.Init(),
		model.signalLoadGraph(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q":
			return model, tea.Quit

		case "down":
			updatedNodeList, command := model.nodeList.TypedUpdate(listview.MsgSelectNext{})
			model.nodeList = updatedNodeList
			return model, command

		case "up":
			updatedNodeList, command := model.nodeList.TypedUpdate(listview.MsgSelectPrevious{})
			model.nodeList = updatedNodeList
			return model, command
		}

	case tea.WindowSizeMsg:
		return model.onResized(message)

	case MsgGraphLoaded:
		return model.onGraphLoaded(message)

	default:
		updatedNodeList, command := model.nodeList.TypedUpdate(message)
		model.nodeList = updatedNodeList
		return model, command
	}

	return model, nil
}

func (model Model) View() string {
	return model.nodeList.View()
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

	updatedNodeList, command := model.nodeList.TypedUpdate(message)
	model.nodeList = updatedNodeList

	return model, command
}

func (model Model) signalUpdateNodeList() tea.Cmd {
	return func() tea.Msg {
		iterator := model.graph.FindNameMatches("")
		names := []string{}

		for iterator.Current() != nil {
			name := iterator.Current().Name
			names = append(names, name)
			iterator.Next()
		}

		return listview.MsgSetItems{
			Items: &SliceAdapter[string]{
				slice: names,
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
