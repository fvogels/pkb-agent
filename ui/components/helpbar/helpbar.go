package helpbar

import (
	"log/slog"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	keyBindings      []key.Binding
	keyStyle         lipgloss.Style
	descriptionStyle lipgloss.Style
}

func New() Model {
	return Model{
		keyBindings:      nil,
		keyStyle:         lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")).Foreground(lipgloss.Color("#FFFFFF")),
		descriptionStyle: lipgloss.NewStyle().Background(lipgloss.Color("#555555")).Foreground(lipgloss.Color("#FFFFFF")),
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case MsgSetKeyBindings:
		return model.setKeyBindings(message)

	default:
		return model, nil
	}
}

func (model Model) View() string {
	slog.Debug("viewing help bar", "n", len(model.keyBindings))
	parts := []string{}

	for _, keyBinding := range model.keyBindings {
		part := lipgloss.JoinHorizontal(
			0,
			model.keyStyle.Render(" "+keyBinding.Help().Key+" "),
			model.descriptionStyle.Render(" "+keyBinding.Help().Desc+" "),
		)

		parts = append(parts, part)
	}

	return lipgloss.JoinHorizontal(0, parts...)
}

func (model Model) setKeyBindings(message MsgSetKeyBindings) (Model, tea.Cmd) {
	model.keyBindings = message.KeyBindings

	return model, nil
}
