package markdownview

import (
	"log/slog"
	"pkb-agent/ui/uid"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	id               int
	source           []byte
	renderedMarkdown string
	size             util.Size
}

func New() Model {
	model := Model{
		source: nil,
		id:     uid.Generate(),
	}

	return model
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

	case MsgSetSource:
		return model.onSetSource(message)

	case msgRenderingDone:
		if message.recipient == model.id {
			return model.onRenderingDone(message)
		}
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	style := lipgloss.NewStyle().MaxWidth(model.size.Width).MaxHeight(model.size.Height)
	return style.Render(model.renderedMarkdown)
}

// onRenderingDone is only called when the recipient matches the current component's id.
func (model Model) onRenderingDone(message msgRenderingDone) (Model, tea.Cmd) {
	model.renderedMarkdown = message.renderedMarkdown
	return model, nil
}

func (model Model) onSetSource(message MsgSetSource) (Model, tea.Cmd) {
	slog.Debug("set source", "id", model.id, "source", message.Source)

	width := model.size.Width
	recipient := model.id

	command := func() tea.Msg {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(width-2),
		)
		if err != nil {
			panic("failed to create markdown renderer")
		}
		renderedMarkdown, err := renderer.Render(message.Source)
		if err != nil {
			panic("failed to render markdown file")
		}
		return msgRenderingDone{
			recipient:        recipient,
			renderedMarkdown: renderedMarkdown,
		}
	}

	return model, command
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}
	return model, nil
}
