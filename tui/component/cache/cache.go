package cache

import (
	"pkb-agent/tui"
)

type Component struct {
	size   tui.Size
	child  tui.Component
	cached tui.Grid
	dirty  bool
}

func New(child tui.Component) *Component {
	return &Component{
		child: child,
	}
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
		component.rerender()
		component.dirty = false
	}

	return component.cached
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
	component.child.Handle(message)
	component.Dirty()
}

func (component *Component) Dirty() {
	component.dirty = true
}

// rerender asks the child component to rerender itself, which overwrites the cache
func (component *Component) rerender() {
	component.cached = tui.FreezeGrid(component.child.Render())
}
