package listview

import (
	"log/slog"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List interface {
	At(index int) string
	Length() int
}

type Model struct {
	items             List
	allowSelection    bool
	firstVisibleIndex int
	selectedIndex     int
	size              util.Size
	filter            func(item string) bool
	emptyListMessage  string
	selectedItemStyle lipgloss.Style
}

func New(allowSelection bool) Model {
	model := Model{
		allowSelection:    allowSelection,
		items:             nil,
		firstVisibleIndex: 0,
		selectedIndex:     0,
		size:              util.Size{Width: 0, Height: 0},
		filter:            func(item string) bool { return true },
		selectedItemStyle: lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")),
		emptyListMessage:  "no nodes found",
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
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case MsgSetItems:
		return model.onSetItems(message)

	case MsgSelectPrevious:
		return model.onSelectPrevious()

	case MsgSelectNext:
		return model.onSelectNext()

	case tea.WindowSizeMsg:
		return model.onResize(message)
	}

	return model, nil
}

func (model Model) View() string {
	if model.items == nil || model.items.Length() == 0 {
		return model.emptyListMessage
	}

	visibleItems := []string{}
	currentIndex := model.firstVisibleIndex
	accumulatedHeight := 0

	for currentIndex < model.items.Length() && accumulatedHeight < model.size.Height {
		item := model.items.At(currentIndex)

		if model.allowSelection && currentIndex == model.selectedIndex {
			item = model.selectedItemStyle.Width(model.size.Width).MaxWidth(model.size.Width).Render(item)
		}

		visibleItems = append(visibleItems, item)

		currentIndex += 1
		accumulatedHeight += 1
	}

	style := lipgloss.NewStyle().Width(model.size.Width).MaxWidth(model.size.Width).Height(model.size.Height).MaxHeight(model.size.Height)
	return style.Render(lipgloss.JoinVertical(0, visibleItems...))
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	model.ensureSelectedIsVisible()

	return model, nil
}

func (model *Model) ensureSelectedIsVisible() {
	if model.allowSelection {
		if model.selectedIndex < model.firstVisibleIndex {
			model.firstVisibleIndex = model.selectedIndex
		} else if model.firstVisibleIndex+model.size.Height <= model.selectedIndex {
			model.firstVisibleIndex = model.selectedIndex - model.size.Height + 1
		}
	}
}

func (model *Model) signalItemSelected() tea.Cmd {
	index := model.selectedIndex

	if model.items.Length() == 0 {
		index = -1
	}

	return func() tea.Msg {
		return MsgItemSelected{
			Index: index,
		}
	}
}

func (model Model) onSetItems(message MsgSetItems) (Model, tea.Cmd) {
	model.items = message.Items
	model.selectedIndex = 0
	model.firstVisibleIndex = 0
	return model, model.signalItemSelected()
}

func (model Model) onSelectPrevious() (Model, tea.Cmd) {
	if model.allowSelection {
		if model.selectedIndex > 0 {
			model.selectedIndex--
		}
		model.ensureSelectedIsVisible()

		return model, model.signalItemSelected()
	} else {
		return model, nil
	}
}

func (model Model) onSelectNext() (Model, tea.Cmd) {
	if model.allowSelection {
		if model.selectedIndex != -1 && model.selectedIndex+1 < model.items.Length() {
			model.selectedIndex++
		}
		model.ensureSelectedIsVisible()

		slog.Debug("new index", slog.Int("index", model.firstVisibleIndex))

		return model, model.signalItemSelected()
	} else {
		return model, nil
	}
}
