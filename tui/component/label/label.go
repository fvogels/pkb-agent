package label

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"

	"github.com/gdamore/tcell/v3"
)

type Component struct {
	tui.ComponentBase
	contents data.Value[string]
	style    *tui.Style
}

func New(messageQueue tui.MessageQueue, name string, contents data.Value[string]) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         name,
			MessageQueue: messageQueue,
		},
		contents: contents,
		style:    &tcell.StyleDefault,
	}

	return &component
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
	return newGrid(component, []rune(component.contents.Get()))
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
}
