package mainscreen

import (
	"pkb-agent/ui/components/nodeselectionview"
	"pkb-agent/ui/layout"
	"pkb-agent/ui/layout/border"
	"pkb-agent/ui/layout/vertical"
	"pkb-agent/ui/nodeviewers/nodeviewer"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var viewModeKeyMap = struct {
	Quit                    key.Binding
	SwitchToInputMode       key.Binding
	HighlightFirst          key.Binding
	HighlightLast           key.Binding
	HighlightNext           key.Binding
	HighlightPrevious       key.Binding
	HighlightPageDown       key.Binding
	HighlightPageUp         key.Binding
	Select                  key.Binding
	UnselectLast            key.Binding
	GrowNodeSelectionView   key.Binding
	ShrinkNodeSelectionView key.Binding
	ToggleIncludedLinked    key.Binding
}{
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	SwitchToInputMode: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
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
		key.WithHelp("↑", "next"),
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
	UnselectLast: key.NewBinding(
		key.WithKeys("delete"),
		key.WithHelp("del", "pop"),
	),
	GrowNodeSelectionView: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "grow node list"),
	),
	ShrinkNodeSelectionView: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "shrink node list"),
	),
	ToggleIncludedLinked: key.NewBinding(
		key.WithKeys("*"),
		key.WithHelp("*", "toggle linked"),
	),
}

type viewMode struct {
	layout layout.Layout[Model]
}

func NewViewMode(layoutConfiguration *layoutConfiguration) *viewMode {
	vlayout := vertical.New[Model]()

	vlayout.Add(
		func(_ util.Size) int { return layoutConfiguration.nodeSelectionViewHeight },
		border.New(layout.Wrap(func(m *Model) *nodeselectionview.Model { return &m.nodeSelectionView })),
	)
	vlayout.Add(
		func(size util.Size) int { return size.Height - layoutConfiguration.nodeSelectionViewHeight },
		border.New(layout.Wrap(func(m *Model) *nodeviewer.Model { return &m.nodeViewer })),
	)

	return &viewMode{
		layout: &vlayout,
	}
}

func (mode viewMode) onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, viewModeKeyMap.Quit):
		return model, tea.Quit

	case key.Matches(message, viewModeKeyMap.SwitchToInputMode):
		model.mode = model.inputMode
		command := model.mode.activate(&model)
		return model, command

	case key.Matches(message, viewModeKeyMap.HighlightFirst):
		return model.onSelectFirstRemainingNode()

	case key.Matches(message, viewModeKeyMap.HighlightLast):
		return model.onSelectLastRemainingNode()

	case key.Matches(message, viewModeKeyMap.HighlightNext):
		return model.onHighlightNextRemainingNode()

	case key.Matches(message, viewModeKeyMap.HighlightPrevious):
		return model.onHighlightPreviousRemainingNode()

	case key.Matches(message, inputModeKeyMap.HighlightPageDown):
		return model.onHighlightRemainingNodePageDown()

	case key.Matches(message, inputModeKeyMap.HighlightPageUp):
		return model.onHighlightRemainingNodePageUp()

	case key.Matches(message, viewModeKeyMap.Select):
		return model.onNodeSelected()

	case key.Matches(message, viewModeKeyMap.UnselectLast):
		return model.onUnselectLast()

	case key.Matches(message, viewModeKeyMap.GrowNodeSelectionView):
		return model.updateLayoutConfiguration(func(c *layoutConfiguration) {
			c.nodeSelectionViewHeight++
		})

	case key.Matches(message, viewModeKeyMap.ToggleIncludedLinked):
		return model.toggleIncludeLinkedNodes()

	case key.Matches(message, viewModeKeyMap.ShrinkNodeSelectionView):
		return model.updateLayoutConfiguration(func(c *layoutConfiguration) {
			if c.nodeSelectionViewHeight > 5 {
				c.nodeSelectionViewHeight--
			}
		})

	default:
		return util.UpdateSingleChild(&model, &model.nodeViewer, message)
	}
}

func (mode viewMode) render(model *Model) string {
	return mode.layout.LayoutView(model)
}

func (mode viewMode) activate(model *Model) tea.Cmd {
	return mode.layout.LayoutResize(model, model.size)
}

func (mode viewMode) resize(model *Model, size util.Size) tea.Cmd {
	return mode.layout.LayoutResize(model, model.size)
}
