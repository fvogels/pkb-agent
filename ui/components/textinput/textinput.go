package textinput

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	input string
}

func New() Model {
	return Model{}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		return model.onKeyPressed(message)

	default:
		return model, nil
	}
}

func (model Model) View() string {
	return model.input
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "backspace":
		if len(model.input) > 0 {
			model.input = model.input[:len(model.input)-1]
			return model, model.signalUpdate()
		} else {
			return model, nil
		}

	default:
		if len(message.String()) == 1 {
			model.input += message.String()
			return model, model.signalUpdate()
		} else {
			return model, nil
		}
	}
}

func (model *Model) GetInput() string {
	return model.input
}

func (model *Model) signalUpdate() tea.Cmd {
	return func() tea.Msg {
		return MsgInputUpdated{
			Input: model.input,
		}
	}
}
