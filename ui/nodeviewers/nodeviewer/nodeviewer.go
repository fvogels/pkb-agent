package nodeviewer

import (
	"log/slog"
	"pkb-agent/graph/nodes/atom"
	"pkb-agent/graph/nodes/bookmark"
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers/bookmarkviewer"
	"pkb-agent/ui/nodeviewers/nullviewer"
	"pkb-agent/ui/nodeviewers/snippetviewer"
	"pkb-agent/util"
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size   util.Size
	viewer tea.Model
}

func New() Model {
	return Model{
		viewer: nullviewer.New(),
	}
}

func (model Model) Init() tea.Cmd {
	return model.viewer.Init()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case MsgSetNode:
		return model.onSetNode(message)

	default:
		return util.UpdateSingleUntypedChild(&model, &model.viewer, message)
	}
}

func (model Model) View() string {
	return model.viewer.View()
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	updatedViewer, command := model.viewer.Update(message)
	model.viewer = updatedViewer

	return model, command
}

func (model Model) onSetNode(message MsgSetNode) (Model, tea.Cmd) {
	node := message.Node

	switch nodeData := node.Extra.(type) {
	case *atom.Extra:
		model.viewer = nullviewer.New()

	case *snippet.Extra:
		model.viewer = snippetviewer.New(nodeData)

	case *bookmark.Extra:
		model.viewer = bookmarkviewer.New(nodeData)

	default:
		slog.Debug(
			"unrecognized node type",
			slog.String("type", reflect.TypeOf(node.Extra).String()),
		)

		model.viewer = nullviewer.New()
	}

	commands := []tea.Cmd{}

	commands = append(commands, model.viewer.Init())
	util.UpdateUntypedChild(&model.viewer, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.size.Height,
	}, &commands)

	return model, tea.Batch(commands...)
}
