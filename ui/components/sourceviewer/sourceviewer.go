package sourceviewer

import (
	"log/slog"
	"pkb-agent/util"
	"pkb-agent/util/syntaxhighlighting"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size            util.Size
	source          string
	language        string
	formattedSource string
}

func New() Model {
	model := Model{
		source:          "",
		formattedSource: "",
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

	case msgSourceFormatted:
		model.formattedSource = message.formattedSource
		return model, nil
	}

	return model, nil
}

func (model Model) View() string {
	style := lipgloss.NewStyle().MaxWidth(model.size.Width).MaxHeight(model.size.Height)

	if len(model.formattedSource) == 0 {
		return style.Render(model.source)
	} else {
		return style.Render(model.formattedSource)
	}
}

func (model Model) onSetSource(message MsgSetSource) (Model, tea.Cmd) {
	model.language = message.Language
	model.source = message.Source

	return model, model.signalFormatSource()
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}
	return model, nil
}

func (model Model) signalFormatSource() tea.Cmd {
	return func() tea.Msg {
		source := model.source
		language := model.language

		formattedSource, err := syntaxhighlighting.Highlight(source, language)
		if err != nil {
			slog.Error("Failed to highlight source code")
			panic("failed to highlight source code")
		}

		return msgSourceFormatted{
			formattedSource: formattedSource,
		}
	}
}
