package holder

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/tui/debug"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child data.Value[tui.Component]
}

func New(messageQueue tui.MessageQueue, child data.Value[tui.Component]) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless holder",
			MessageQueue: messageQueue,
		},
		child: child,
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	child := component.child.Get()

	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgStateUpdated:
		if child != nil {
			child.Handle(message)
			child.Handle(tui.MsgResize{Size: component.Size})
		}

	default:
		if child != nil {
			component.child.Get().Handle(message)
		}
	}
}

func (component *Component) Render() grid.Grid {
	child := component.child.Get()

	if child != nil {
		return child.Render()
	} else {
		return tui.NewEmptyGrid(component.Size)
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	child := component.child.Get()
	if child != nil {
		child.Handle(message)
	}
}
