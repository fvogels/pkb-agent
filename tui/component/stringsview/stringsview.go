package stringsview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/tui/debug"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/size"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	items             data.Value[list.List[Item]]
	emptyStyle        *tui.Style
	firstVisibleIndex int
	onItemClicked     func(int)
}

type Item struct {
	Runes []rune
	Style *tui.Style
}

func New(messageQueue tui.MessageQueue, items data.Value[list.List[Item]]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed strings view",
			MessageQueue: messageQueue,
		},
		items:             items,
		emptyStyle:        &defaultEmptyStyle,
		firstVisibleIndex: 0,
		onItemClicked:     nil,
	}

	return &component
}

func (component *Component) SetEmptyStyle(emptyStyle *tui.Style) {
	component.emptyStyle = emptyStyle
}

func (component *Component) SetOnItemClicked(onItemClicked func(int)) {
	component.onItemClicked = onItemClicked
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tuigrid.Grid {
	return newGrid(component)
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
}

func (component *Component) EnsureItemIsVisible(index int) {
	if component.items.Get().Size() == 0 {
		component.firstVisibleIndex = 0
	} else {
		if index < component.firstVisibleIndex {
			component.firstVisibleIndex = index
		} else if component.firstVisibleIndex+component.Size.Height <= index {
			component.firstVisibleIndex = index - component.Size.Height + 1
		}
	}
}

func (component *Component) Measure() size.Size {
	measuredWidth := 0
	items := component.items.Get()

	for index := range items.Size() {
		item := items.At(index)
		itemWidth := len(item.Runes)

		if itemWidth > measuredWidth {
			measuredWidth = itemWidth
		}
	}

	return size.Size{
		Width:  measuredWidth,
		Height: items.Size(),
	}
}
