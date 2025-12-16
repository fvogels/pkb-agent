package mainscreen

import (
	"pkb-agent/ui/components/nodeselectionview"
	"pkb-agent/ui/components/textinput"
	"pkb-agent/ui/layout"
	"pkb-agent/ui/layout/border"
	"pkb-agent/ui/layout/vertical"
	"pkb-agent/ui/nodeviewers/nodeviewer"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var inputModeKeyMap = struct {
	Cancel            key.Binding
	HighlightNext     key.Binding
	HighlightPrevious key.Binding
	Select            key.Binding
}{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	HighlightNext: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next"),
	),
	HighlightPrevious: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "next"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↵", "select"),
	),
}

type inputMode struct {
	layout layout.Layout[Model]
}

func NewInputMode(layoutConfiguration *layoutConfiguration) *inputMode {
	vlayout := vertical.New[Model]()

	vlayout.Add(
		func(_ util.Size) int { return layoutConfiguration.nodeSelectionViewHeight },
		border.New(layout.Wrap(func(m *Model) *nodeselectionview.Model { return &m.nodeSelectionView })),
	)
	vlayout.Add(
		func(size util.Size) int { return size.Height - layoutConfiguration.nodeSelectionViewHeight - 1 },
		border.New(layout.Wrap(func(m *Model) *nodeviewer.Model { return &m.nodeViewer })),
	)
	vlayout.Add(
		func(_ util.Size) int { return 1 },
		layout.Wrap(func(m *Model) *textinput.Model { return &m.textInput }),
	)

	return &inputMode{
		layout: &vlayout,
	}
}

func (mode inputMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, inputModeKeyMap.Cancel):
		var command tea.Cmd
		model, command = util.UpdateSingleChild(&model, &model.textInput, textinput.MsgClear{})

		model.mode = model.viewMode
		model.viewMode.activate(&model)
		return model, command

	case key.Matches(message, inputModeKeyMap.HighlightNext):
		return model.onSelectNextRemainingNode()

	case key.Matches(message, inputModeKeyMap.HighlightPrevious):
		return model.onSelectPreviousRemainingNode()

	case key.Matches(message, inputModeKeyMap.Select):
		model.mode = model.viewMode
		model.mode.activate(&model)
		return model.onNodeSelected()

	default:
		updatedTextInput, command := model.textInput.TypedUpdate(message)
		model.textInput = updatedTextInput

		return model, command
	}
}

func (mode inputMode) render(model *Model) string {
	return mode.layout.LayoutView(model)
}

func (mode inputMode) activate(model *Model) tea.Cmd {
	return mode.layout.LayoutResize(model, model.size)
}

func (mode inputMode) resize(model *Model, size util.Size) tea.Cmd {
	return mode.layout.LayoutResize(model, model.size)
}
