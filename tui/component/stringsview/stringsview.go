package stringsview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	size              tui.Size
	items             data.Value[list.List[Item]]
	emptyStyle        *tui.Style
	firstVisibleIndex int
	onItemClicked     func(int)
}

type Item struct {
	Runes []rune
	Style *tui.Style
}

func New(items data.Value[list.List[Item]]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	return &Component{
		items:             items,
		emptyStyle:        &defaultEmptyStyle,
		firstVisibleIndex: 0,
		onItemClicked:     nil,
	}
}

func (component *Component) SetEmptyStyle(emptyStyle *tui.Style) {
	component.emptyStyle = emptyStyle
}

func (component *Component) SetOnItemClicked(onItemClicked func(int)) {
	component.onItemClicked = onItemClicked
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return newGrid(component)
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
}

func (component *Component) EnsureItemIsVisible(index int) {
	if component.items.Get().Size() == 0 {
		component.firstVisibleIndex = 0
	} else {
		if index < component.firstVisibleIndex {
			component.firstVisibleIndex = index
		} else if component.firstVisibleIndex+component.size.Height <= index {
			component.firstVisibleIndex = index - component.size.Height + 1
		}
	}
}
