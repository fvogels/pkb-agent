package stringlist

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/component/stringsview"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	items              data.Value[list.List[string]]
	selectedIndex      data.Value[int]
	emptyStyle         *tui.Style
	itemStyle          *tui.Style
	selectedItemStyle  *tui.Style
	firstVisibleIndex  int
	onSelectionChanged func(int)
	child              *stringsview.Component
}

func New(messageQueue tui.MessageQueue, items data.Value[list.List[string]], selectedItem data.Value[int]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultItemStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	defaultSelectedItemStyle := tcell.StyleDefault.Background(color.Gray).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed stringlist",
			MessageQueue: messageQueue,
		},
		items:             items,
		selectedIndex:     selectedItem,
		itemStyle:         &defaultItemStyle,
		emptyStyle:        &defaultEmptyStyle,
		selectedItemStyle: &defaultSelectedItemStyle,
		firstVisibleIndex: 0,
	}

	subComponentList := data.MapValue2(items, selectedItem, func(items list.List[string], selectedIndex int) list.List[stringsview.Item] {
		return list.MapWithIndex(items, func(index int, item string) stringsview.Item {
			var style *tui.Style
			if index == selectedIndex {
				style = component.selectedItemStyle
			} else {
				style = component.itemStyle
			}

			return stringsview.Item{
				Runes: []rune(item),
				Style: style,
			}
		})
	})

	component.child = stringsview.New(messageQueue, subComponentList)

	// update selection when item clicked
	component.child.SetOnItemClicked(func(index int) {
		if component.onSelectionChanged != nil {
			component.onSelectionChanged(index)
		}
	})

	// ensure that selected item is visible at all times
	selectedItem.Observe(func() {
		component.ensureSelectedItemIsVisible()
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
		component.child.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
	component.ensureSelectedItemIsVisible()
}

func (component *Component) ensureSelectedItemIsVisible() {
	selectedIndex := component.selectedIndex.Get()

	component.child.EnsureItemIsVisible(selectedIndex)
}

func (component *Component) onKey(message tui.MsgKey) {
	selectedIndex := component.selectedIndex.Get()
	maximumIndex := component.items.Get().Size() - 1
	pageSize := component.Size.Height
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
}

type SubComponentList struct {
	items         data.Value[list.List[string]]
	selectedIndex data.Value[int]
	defaultStyle  *tui.Style
	selectedStyle *tui.Style
}

func (list *SubComponentList) Size() int {
	return list.items.Get().Size()
}

func (list *SubComponentList) At(index int) stringsview.Item {
	var style *tui.Style
	if index == list.selectedIndex.Get() {
		style = list.selectedStyle
	} else {
		style = list.defaultStyle
	}

	return stringsview.Item{
		Runes: []rune(list.items.Get().At(index)),
		Style: style,
	}
}
