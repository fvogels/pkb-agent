package mainscreen

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func New() Model {
	return Model{}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.KeyMsg:
		switch message.String() {
		case "q":
			return model, tea.Quit
		}
	}

	return model, nil
}

func (model Model) View() string {
	return "hello world"
}
