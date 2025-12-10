package nodeviewer

import (
	"pkb-agent/ui/nodeviewers/nullviewer"
	"pkb-agent/util"

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
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case MsgSetNode:
		return model.onSetNode(message)
	}

	return model, nil
}

func (model Model) View() string {
	return model.viewer.View()
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	updatedViewer, command := model.Update(message)
	model.viewer = updatedViewer

	return model, command
}

func (model Model) onSetNode(message MsgSetNode) (Model, tea.Cmd) {
	return model, nil
	// node := message.Node

	// switch nodeData := node.Extra.(type) {
	// case *atom.Extra:
	// 	model.viewer = nullviewer.New()

	// case *snippet.Extra:
	// 	model.viewer = snippetviewer.New(nodeData)

	// default:
	// 	model.viewer = nullviewer.New()
	// }

	// commands := []tea.Cmd{}

	// commands = append(commands, model.viewer.Init())
	// util.UpdateUntypedChild(&model.viewer, tea.WindowSizeMsg{
	// 	Width:  model.size.Width,
	// 	Height: model.size.Height,
	// }, &commands)

	// return model, tea.Batch(commands...)
}
