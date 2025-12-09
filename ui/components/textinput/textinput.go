package textinput

import (
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size   util.Size
	input  string
	style  lipgloss.Style
	cursor cursor.Model
}

func New() Model {
	c := cursor.New()
	c.SetChar("█")

	return Model{
		style:  lipgloss.NewStyle().Background(lipgloss.Color("#5555FF")),
		cursor: c,
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.cursor.Focus(),
		model.cursor.SetMode(cursor.CursorStatic),
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

	case MsgClear:
		model.input = ""
		return model, model.signalUpdate()

	case tea.WindowSizeMsg:
		return model.onResize(message)

	default:
		updatedCursor, command := model.cursor.Update(message)
		model.cursor = updatedCursor
		return model, command
	}
}

func (model Model) View() string {
	return model.style.Width(model.size.Width).Render(
		lipgloss.JoinHorizontal(
			0,
			model.input,
			model.cursor.View(),
		),
	)
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

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}
