package holder

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
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

	// Make sure that whenever a new component is put in, it is activated and resized
	child.Observe(func() {
		component.child.Get().Handle(tui.MsgActivate{})
		component.MessageQueue.Enqueue(tui.MsgUpdateLayout{})
	})

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.child.Get().Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
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
