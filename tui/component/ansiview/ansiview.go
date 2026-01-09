package ansiview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
)

type Component struct {
	size        tui.Size
	rawContents data.Value[string]
	ansiGrid    data.Value[tui.Grid]
}

func New(contents data.Value[string]) *Component {
	return &Component{
		rawContents: contents,
	}
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return newGrid(component.size, component.ansiGrid.Get())
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
}
