package label

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
)

type Component struct {
	size     tui.Size
	contents data.Data[string]
	style    tui.Style
}

func New(contents data.Data[string], style tui.Style) *Component {
	return &Component{
		contents: contents,
		style:    style,
	}
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
		style:    &component.style,
		size:     component.size,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
}
