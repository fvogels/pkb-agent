package stringsview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	size              tui.Size
	items             data.List[Item]
	emptyStyle        *tui.Style
	firstVisibleIndex int
}

type Item struct {
	Runes []rune
	Style *tui.Style
}

func New(items data.List[Item]) *Component {
	defaultEmptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	return &Component{
		items:             items,
		emptyStyle:        &defaultEmptyStyle,
		firstVisibleIndex: 0,
	}
}

func (component *Component) SetEmptyStyle(emptyStyle *tui.Style) {
	component.emptyStyle = emptyStyle
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
