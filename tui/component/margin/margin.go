package margin

import (
	"pkb-agent/tui"
	tuigrid "pkb-agent/tui/grid"
	"pkb-agent/tui/size"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child      tui.Component
	marginSize int
	style      *tui.Style
}

func New(messageQueue tui.MessageQueue, child tui.Component, marginSize int, style *tui.Style) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed border",
			MessageQueue: messageQueue,
		},
		child:      child,
		marginSize: marginSize,
		style:      style,
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.child.Handle(message)
	}
}

func (component *Component) Render() tuigrid.FiniteGrid {
	return newGrid(component)
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	childSizeMessage := tui.MsgResize{
		Size: size.Size{
			Width:  message.Size.Width - 2,
			Height: component.Size.Height - 2,
		},
	}
	component.child.Handle(childSizeMessage)
}
