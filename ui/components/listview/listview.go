package listview

import (
	"fmt"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type List[T any] interface {
	At(index int) T
	Length() int
}

type Model[T any] struct {
	itemRenderer           func(item T) string
	items                  List[T]
	allowSelection         bool
	firstVisibleIndex      int
	selectedIndex          int
	size                   util.Size
	emptyListMessage       string
	nonselectedItemStyle   lipgloss.Style
	selectedItemStyle      lipgloss.Style
	outgoingMessageWrapper func(tea.Msg) tea.Msg
}

func New[T any](itemRenderer func(item T) string, allowSelection bool, outgoingMessageWrapper func(tea.Msg) tea.Msg) Model[T] {
	model := Model[T]{
		itemRenderer:           itemRenderer,
		allowSelection:         allowSelection,
		items:                  &emptyList[T]{},
		firstVisibleIndex:      0,
		selectedIndex:          0,
		size:                   util.Size{Width: 0, Height: 0},
		nonselectedItemStyle:   lipgloss.NewStyle(),
		selectedItemStyle:      lipgloss.NewStyle().Background(lipgloss.Color("#CCCCCC")).Foreground(lipgloss.Color("#000000")),
		emptyListMessage:       "no nodes found",
		outgoingMessageWrapper: outgoingMessageWrapper,
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
	if model.items.Length() == 0 {
		return model.emptyListMessage
	}

	visibleItems := []string{}
	currentIndex := model.firstVisibleIndex
	accumulatedHeight := 0

	selectedStyle := model.selectedItemStyle.Width(model.size.Width).MaxWidth(model.size.Width)
	nonselectedStyle := model.nonselectedItemStyle.Width(model.size.Width).MaxWidth(model.size.Width)

	for currentIndex < model.items.Length() && accumulatedHeight < model.size.Height {
		item := model.itemRenderer(model.items.At(currentIndex))

		var style lipgloss.Style
		if model.allowSelection && currentIndex == model.selectedIndex {
			style = selectedStyle
		} else {
			style = nonselectedStyle
		}

		item = style.Render(item)

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
	if model.allowSelection && model.items != nil {
		if model.selectedIndex < model.firstVisibleIndex {
			model.firstVisibleIndex = model.selectedIndex
		} else if model.firstVisibleIndex+model.size.Height <= model.selectedIndex {
			model.firstVisibleIndex = model.selectedIndex - model.size.Height + 1
		}
	}
}

func (model *Model[T]) signalItemSelected() tea.Cmd {
	index := model.selectedIndex

	// Sanity check
	if index >= model.items.Length() {
		panic(fmt.Sprintf("invalid index %d in listview of size %d", index, model.items.Length()))
	}

	selectedItem := model.items.At(index)

	return func() tea.Msg {
		message := MsgItemSelected[T]{
			Index: index,
			Item:  selectedItem,
		}

		return model.outgoingMessageWrapper(message)
	}
}

func (model *Model[T]) signalNoItemsSelected() tea.Cmd {
	return func() tea.Msg {
		message := MsgNoItemSelected{}
		return model.outgoingMessageWrapper(message)
	}
}

func (model Model[T]) onSetItems(message MsgSetItems[T]) (Model[T], tea.Cmd) {
	model.items = message.Items

	if !model.allowSelection {
		return model, nil
	}

	model.selectedIndex = message.SelectionIndex
	model.firstVisibleIndex = util.MaxInt(0, message.SelectionIndex-2)
	model.ensureSelectedIsVisible()

	if model.items.Length() > 0 {
		return model, model.signalItemSelected()
	} else {
		return model, model.signalNoItemsSelected()
	}
}

func (model Model[T]) onSelectPrevious() (Model[T], tea.Cmd) {
	if model.allowSelection && model.items.Length() > 0 {
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
	if model.allowSelection && model.items.Length() > 0 {
		if model.selectedIndex+1 < model.items.Length() {
			model.selectedIndex++
		}
		model.ensureSelectedIsVisible()

		return model, model.signalItemSelected()
	} else {
		return model, nil
	}
}

func (model *Model[T]) GetSelectedIndex() int {
	return model.selectedIndex
}

func (model *Model[T]) GetSelectedItem() T {
	return model.items.At(model.selectedIndex)
}

func (model *Model[T]) SetSelectedStyle(style lipgloss.Style) {
	model.selectedItemStyle = style
}

func (model *Model[T]) SetNonselectedStyle(style lipgloss.Style) {
	model.nonselectedItemStyle = style
}

func (model *Model[T]) SetEmptyListMessage(message string) {
	model.emptyListMessage = message
}
