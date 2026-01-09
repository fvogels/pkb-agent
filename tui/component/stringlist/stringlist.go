package stringlist

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/stringsview"
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
	subComponent       *stringsview.Component
}

func New(items data.List[string], selectedItem data.Value[int]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultItemStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultSelectedItemStyle := tcell.StyleDefault.Background(color.Gray).Foreground(color.Reset)
	subComponentList := SubComponentList{
		items:         items,
		selectedIndex: selectedItem,
		defaultStyle:  &defaultItemStyle,
		selectedStyle: &defaultSelectedItemStyle,
	}
	subComponent := stringsview.New(&subComponentList)

	component := Component{
		items:             items,
		selectedIndex:     selectedItem,
		itemStyle:         &defaultItemStyle,
		emptyStyle:        &defaultEmptyStyle,
		selectedItemStyle: &defaultSelectedItemStyle,
		firstVisibleIndex: 0,
		subComponent:      subComponent,
	}

	subComponent.SetOnItemClicked(func(index int) {
		if component.onSelectionChanged != nil {
			component.onSelectionChanged(index)
		}
	})

	return &component
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

	default:
		component.subComponent.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	return component.subComponent.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
	component.subComponent.Handle(message)
	component.ensureSelectedItemIsVisible()
}

func (component *Component) ensureSelectedItemIsVisible() {
	selectedIndex := component.selectedIndex.Get()

	component.subComponent.EnsureItemIsVisible(selectedIndex)
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

		if component.onSelectionChanged != nil {
			component.onSelectionChanged(index)
		}
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

	component.ensureSelectedItemIsVisible()
}

type SubComponentList struct {
	items         data.List[string]
	selectedIndex data.Value[int]
	defaultStyle  *tui.Style
	selectedStyle *tui.Style
}

func (list *SubComponentList) Size() int {
	return list.items.Size()
}

func (list *SubComponentList) At(index int) stringsview.Item {
	var style *tui.Style
	if index == list.selectedIndex.Get() {
		style = list.selectedStyle
	} else {
		style = list.defaultStyle
	}

	return stringsview.Item{
		Runes: []rune(list.items.At(index)),
		Style: style,
	}
}

func (list *SubComponentList) Observe(func()) {}
