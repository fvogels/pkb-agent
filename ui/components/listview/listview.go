package listview

import (
	"log/slog"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List[T Item] interface {
	At(index int) T
	Length() int
}

type Item interface {
	String() string
}

type Model[T Item] struct {
	items             List[T]
	allowSelection    bool
	firstVisibleIndex int
	selectedIndex     int
	size              util.Size
	emptyListMessage  string
	selectedItemStyle lipgloss.Style
}

func New[T Item](allowSelection bool) Model[T] {
	model := Model[T]{
		allowSelection:    allowSelection,
		items:             nil,
		firstVisibleIndex: 0,
		selectedIndex:     0,
		size:              util.Size{Width: 0, Height: 0},
		selectedItemStyle: lipgloss.NewStyle().Background(lipgloss.Color("#AAAAAA")),
		emptyListMessage:  "no nodes found",
	}

	return model
}

func (model Model[T]) Init() tea.Cmd {
	return nil
}

func (model Model[T]) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model[T]) TypedUpdate(message tea.Msg) (Model[T], tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case MsgSetItems[T]:
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

func (model Model[T]) View() string {
	if model.items == nil || model.items.Length() == 0 {
		return model.emptyListMessage
	}

	visibleItems := []string{}
	currentIndex := model.firstVisibleIndex
	accumulatedHeight := 0

	for currentIndex < model.items.Length() && accumulatedHeight < model.size.Height {
		item := model.items.At(currentIndex).String()

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

func (model Model[T]) onResize(message tea.WindowSizeMsg) (Model[T], tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	model.ensureSelectedIsVisible()

	return model, nil
}

func (model *Model[T]) ensureSelectedIsVisible() {
	if model.allowSelection {
		if model.selectedIndex < model.firstVisibleIndex {
			model.firstVisibleIndex = model.selectedIndex
		} else if model.firstVisibleIndex+model.size.Height <= model.selectedIndex {
			model.firstVisibleIndex = model.selectedIndex - model.size.Height + 1
		}
	}
}

func (model *Model[T]) signalItemSelected() tea.Cmd {
	index := model.selectedIndex
	var selectedItem T

	if model.items.Length() == 0 {
		index = -1
	} else {
		selectedItem = model.items.At(index)
	}

	return func() tea.Msg {
		return MsgItemSelected[T]{
			Index: index,
			Item:  selectedItem,
		}
	}
}

func (model Model[T]) onSetItems(message MsgSetItems[T]) (Model[T], tea.Cmd) {
	model.items = message.Items
	model.selectedIndex = 0
	model.firstVisibleIndex = 0
	return model, model.signalItemSelected()
}

func (model Model[T]) onSelectPrevious() (Model[T], tea.Cmd) {
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

func (model Model[T]) onSelectNext() (Model[T], tea.Cmd) {
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

func (model Model[T]) GetSelectedIndex() int {
	return model.selectedIndex
}
