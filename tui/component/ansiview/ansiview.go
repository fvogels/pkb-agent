package ansiview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/ansigrid"
	"pkb-agent/tui/data"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	rawContents data.Value[string]
	ansiGrid    data.Value[tuigrid.Grid]
	emptyStyle  *tui.Style
}

func New(messageQueue tui.MessageQueue, contents data.Value[string]) *Component {
	emptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			MessageQueue: messageQueue,
			Name:         "nameless ansiview",
		},
		rawContents: contents,
		ansiGrid: data.MapValue(contents, func(s string) tuigrid.Grid {
			return ansigrid.Parse(s, &emptyStyle)
		}),
		emptyStyle: &emptyStyle,
	}

	return &component
}

func (component *Component) SetEmptyStyle(style *tui.Style) {
	component.emptyStyle = style
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tuigrid.Grid {
	size := component.Size
	grid := component.ansiGrid.Get()
	style := component.emptyStyle

	return newGrid(size, grid, style)
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
}
