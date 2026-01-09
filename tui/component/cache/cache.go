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

type Observable interface {
	Observe(func())
}

func New(child tui.Component, observables ...Observable) *Component {
	result := Component{
		child: child,
	}

	result.AddInvalidators(observables...)

	return &result
}

func (component *Component) AddInvalidators(observables ...Observable) {
	f := func() { component.Invalidate() }

	for _, observable := range observables {
		observable.Observe(f)
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
	component.Invalidate()
}

func (component *Component) Invalidate() {
	component.dirty = true
}

// rerender asks the child component to rerender itself, which overwrites the cache
func (component *Component) rerender() {
	component.cached = tui.MaterializeGrid(component.child.Render())
}
