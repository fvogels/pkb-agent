package label

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
)

type Component struct {
	name     string
	size     tui.Size
	contents data.Value[string]
	style    *tui.Style
}

func New(name string, contents data.Value[string]) *Component {
	return &Component{
		name:     name,
		contents: contents,
		style:    &tcell.StyleDefault,
	}
}

func (component *Component) SetStyle(style *tui.Style) {
	component.style = style
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return &grid{
		contents: []rune(component.contents.Get()),
		style:    component.style,
		size:     component.size,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
}
