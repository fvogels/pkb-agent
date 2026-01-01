package markdownview

import (
	"log/slog"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/uid"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	id               int
	source           string
	renderedMarkdown string
	size             util.Size
}

func New(source string) Model {
	model := Model{
		id:     uid.Generate(),
		source: source,
	}

	return model
}

func (model Model) Init() tea.Cmd {
	return model.signalFormatMarkdown()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

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
	model.source = message.Source

	return model, model.signalFormatMarkdown()
}

func (model Model) signalFormatMarkdown() tea.Cmd {
	width := model.size.Width
	recipient := model.id
	source := model.source

	command := func() tea.Msg {
		markdownWidth := width - 2

		slog.Debug("formatting markdown", slog.Int("componentWidth", width), slog.Int("markdownWidth", markdownWidth))

		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(markdownWidth),
		)
		if err != nil {
			panic("failed to create markdown renderer")
		}
		renderedMarkdown, err := renderer.Render(source)
		if err != nil {
			panic("failed to render markdown file")
		}
		return msgRenderingDone{
			recipient:        recipient,
			renderedMarkdown: renderedMarkdown,
		}
	}

	return command
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	slog.Debug("!!!")

	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, model.signalFormatMarkdown()
}
