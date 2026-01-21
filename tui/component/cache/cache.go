package cache

import (
	"pkb-agent/tui"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child  tui.Component
	cached tui.Grid
	dirty  bool
}

func New(messageQueue tui.MessageQueue, child tui.Component) *Component {
	result := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless cache",
			MessageQueue: messageQueue,
		},
		child: child,
	}

	return &result
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
	if component.dirty {
		component.Refresh()
		component.dirty = false
	}

	return component.cached
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
	component.Invalidate()
}

func (component *Component) Invalidate() {
	component.dirty = true
}

func (component *Component) Refresh() {
	component.cached = tui.MaterializeGrid(component.child.Render())
}
