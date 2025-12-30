package mainscreen

import (
	"pkb-agent/ui/components/nodeselectionview"
	"pkb-agent/ui/components/nodeviewer"
	"pkb-agent/ui/components/textinput"
	"pkb-agent/ui/layout"
	"pkb-agent/ui/layout/border"
	"pkb-agent/ui/layout/vertical"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var inputModeKeyMap = struct {
	Cancel            key.Binding
	HighlightFirst    key.Binding
	HighlightLast     key.Binding
	HighlightPrevious key.Binding
	HighlightNext     key.Binding
	HighlightPageDown key.Binding
	HighlightPageUp   key.Binding
	Select            key.Binding
}{
	Cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	HighlightFirst: key.NewBinding(
		key.WithKeys("home"),
		key.WithHelp("home", "first"),
	),
	HighlightLast: key.NewBinding(
		key.WithKeys("end"),
		key.WithHelp("end", "last"),
	),
	HighlightNext: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next"),
	),
	HighlightPrevious: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "previous"),
	),
	HighlightPageDown: key.NewBinding(
		key.WithKeys("pgdown"),
		key.WithHelp("pgdown", "pgdown"),
	),
	HighlightPageUp: key.NewBinding(
		key.WithKeys("pgup"),
		key.WithHelp("pgup", "pgup"),
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
		var command1 tea.Cmd
		model, command1 = util.UpdateSingleChild(&model, &model.textInput, textinput.MsgClear{})

		command2 := func() tea.Msg { return msgSwitchMode{mode: model.viewMode} }
		return model, tea.Batch(command1, command2)

	case key.Matches(message, viewModeKeyMap.HighlightFirst):
		return model.onSelectFirstRemainingNode()

	case key.Matches(message, viewModeKeyMap.HighlightLast):
		return model.onSelectLastRemainingNode()

	case key.Matches(message, inputModeKeyMap.HighlightNext):
		return model.onHighlightNextRemainingNode()

	case key.Matches(message, inputModeKeyMap.HighlightPrevious):
		return model.onHighlightPreviousRemainingNode()

	case key.Matches(message, inputModeKeyMap.HighlightPageDown):
		return model.onHighlightRemainingNodePageDown()

	case key.Matches(message, inputModeKeyMap.HighlightPageUp):
		return model.onHighlightRemainingNodePageUp()

	case key.Matches(message, inputModeKeyMap.Select):
		model.mode = model.viewMode
		model.mode.activate(&model)
		return model.onNodeSelected()

	default:
		commands := []tea.Cmd{}

		util.UpdateChild(&model.textInput, message, &commands)
		// util.UpdateChild(&model.nodeViewer, message, &commands)

		return model, tea.Batch(commands...)
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

func (mode inputMode) getKeyBindings() []key.Binding {
	return []key.Binding{}
}
