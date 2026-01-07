package stringlist

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	size               tui.Size
	items              data.List[string]
	selectedIndex      data.Value[int]
	emptyStyle         *tui.Style
	itemStyle          *tui.Style
	selectedItemStyle  *tui.Style
	firstVisibleIndex  int
	onSelectionChanged func(int)
}

func New(items data.List[string], selectedItem data.Value[int]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultItemStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultSelectedItemStyle := tcell.StyleDefault.Background(color.Gray).Foreground(color.Reset)

	return &Component{
		items:             items,
		selectedIndex:     selectedItem,
		itemStyle:         &defaultItemStyle,
		emptyStyle:        &defaultEmptyStyle,
		selectedItemStyle: &defaultSelectedItemStyle,
		firstVisibleIndex: 0,
	}
}

func (component *Component) SetEmptyStyle(emptyStyle *tui.Style) {
	component.emptyStyle = emptyStyle
}

func (component *Component) SetItemStyle(itemStyle *tui.Style) {
	component.itemStyle = itemStyle
}

func (component *Component) SetSelectedItemStyle(selectedItemStyle *tui.Style) {
	component.selectedItemStyle = selectedItemStyle
}

func (component *Component) SetOnSelectionChanged(callback func(int)) {
	component.onSelectionChanged = callback
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgKey:
		component.onKey(message)
	}
}

func (component *Component) Render() tui.Grid {
	return newGrid(component)
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size

	component.ensureSelectedItemIsVisible()
}

func (component *Component) ensureSelectedItemIsVisible() {
	selectedIndex := component.selectedIndex.Get()

	if component.items.Size() == 0 {
		component.firstVisibleIndex = 0
	} else {
		if component.selectedIndex.Get() < component.firstVisibleIndex {
			component.firstVisibleIndex = selectedIndex
		} else if component.firstVisibleIndex+component.size.Height <= selectedIndex {
			component.firstVisibleIndex = selectedIndex - component.size.Height + 1
		}
	}
}

func (component *Component) onKey(message tui.MsgKey) {
	selectedIndex := component.selectedIndex.Get()
	maximumIndex := component.items.Size() - 1
	pageSize := component.size.Height
	onSelectionChanged := func(index int) {
		if index > maximumIndex {
			index = maximumIndex
		}
		if index < 0 {
			index = 0
		}
		component.onSelectionChanged(index)
	}

	switch message.Key {
	case "Down":
		onSelectionChanged(selectedIndex + 1)

	case "Up":
		onSelectionChanged(selectedIndex - 1)

	case "Home":
		onSelectionChanged(0)

	case "End":
		onSelectionChanged(maximumIndex)

	case "PgDn":
		onSelectionChanged(selectedIndex + pageSize)

	case "PgUp":
		onSelectionChanged(selectedIndex - pageSize)
	}
}
