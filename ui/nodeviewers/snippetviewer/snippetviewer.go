package snippetviewer

import (
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

type Model struct {
	size            util.Size
	nodeData        *snippet.Extra
	rawSource       string
	formattedSource string
}

func New(nodeData *snippet.Extra) Model {
	return Model{
		nodeData:        nodeData,
		formattedSource: "loading",
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadSnippet()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case msgSnippetLoaded:
		return model.onSnippetLoaded(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)
	}

	return model, nil
}

func (model Model) View() string {
	return model.formattedSource
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}

func (model *Model) signalLoadSnippet() tea.Cmd {
	data := model.nodeData

	return func() tea.Msg {
		rawSource, formattedSource, err := data.GetHighlightedSource()
		if err != nil {
			panic("failed to get snippet source")
		}

		return msgSnippetLoaded{
			rawSource:       rawSource,
			formattedSource: formattedSource,
		}
	}
}

func (model Model) onSnippetLoaded(message msgSnippetLoaded) (Model, tea.Cmd) {
	model.rawSource = message.rawSource
	model.formattedSource = message.formattedSource
	return model, nil
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch message.String() {
	case "c":
		return model.onCopyToClipboard()

	default:
		return model, nil
	}
}

func (model Model) onCopyToClipboard() (Model, tea.Cmd) {
	clipboard.Write(clipboard.FmtText, []byte(model.rawSource))

	return model, nil
}
