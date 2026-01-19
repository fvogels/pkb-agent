package cache

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child  tui.Component
	cached tui.Grid
	dirty  bool
}

func New(messageQueue tui.MessageQueue, child tui.Component, observables ...data.Observable) *Component {
	result := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless cache",
			MessageQueue: messageQueue,
		},
		child: child,
	}

	result.AddInvalidators(observables...)

	return &result
}

func (component *Component) AddInvalidators(observables ...data.Observable) {
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
	component.Size = message.Size
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
