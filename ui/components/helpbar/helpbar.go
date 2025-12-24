package helpbar

import (
	"fmt"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size               util.Size
	keyBindings        []key.Binding
	remainingNodeCount int
	totalNodeCount     int
	keyStyle           lipgloss.Style
	descriptionStyle   lipgloss.Style
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
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case MsgSetKeyBindings:
		return model.setKeyBindings(message)

	case MsgSetNodeCounts:
		model.totalNodeCount = message.Total
		model.remainingNodeCount = message.Remaining
		return model, nil

	default:
		return model, nil
	}
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}

func (model Model) View() string {
	keyBindings := model.renderKeyBindings()
	counts := model.renderNodeCounts()

	return lipgloss.JoinHorizontal(
		0,
		lipgloss.NewStyle().Width(model.size.Width-lipgloss.Width(counts)).Render(keyBindings),
		counts,
	)
}

func (model Model) renderKeyBindings() string {
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

func (model Model) renderNodeCounts() string {
	return fmt.Sprintf(" %d/%d ", model.remainingNodeCount, model.totalNodeCount)
}

func (model Model) setKeyBindings(message MsgSetKeyBindings) (Model, tea.Cmd) {
	model.keyBindings = message.KeyBindings

	return model, nil
}
