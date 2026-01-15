package border

import (
	"pkb-agent/tui"
	"pkb-agent/ui/uid"
)

type Component struct {
	tui.ComponentBase
	child tui.Component
	style tui.Style
}

func New(messageQueue tui.MessageQueue, child tui.Component, style tui.Style) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed border",
			MessageQueue: messageQueue,
		},
		child: child,
		style: style,
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

func (component *Component) Render() tui.Grid {
	return &grid{
		childGrid: component.child.Render(),
		style:     component.style,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	childSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  message.Size.Width - 2,
			Height: component.Size.Height - 2,
		},
	}
	component.child.Handle(childSizeMessage)
}
