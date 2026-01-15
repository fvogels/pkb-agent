package markdownpage

import (
	"pkb-agent/ui/components/markdownview"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	size           util.Size
	source         string
	markdownViewer markdownview.Model
}

func NewModel(source string) Model {
	return Model{
		source:         source,
		markdownViewer: markdownview.New(source),
	}
}

func (model Model) Init() tea.Cmd {
	return model.markdownViewer.Init()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

	default:
		return util.UpdateSingleChild(&model, &model.markdownViewer, message)
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return util.UpdateSingleChild(&model, &model.markdownViewer, message)
}

func (model Model) onKeyPressed(tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	default:
		return model, nil
	}
}

func (model Model) View() string {
	return model.markdownViewer.View()
}

func (model Model) signalKeybindingsUpdate() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}
