package atom

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child tui.Component
}

func NewViewer(messageQueue tui.MessageQueue) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed atom viewer",
			MessageQueue: messageQueue,
		},
		child: label.New(messageQueue, "atom label", data.NewConstant("atom!")),
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgActivate:
		message.Respond(
			component.Identifier,
			component.onActivate,
			component.child,
		)

	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) onActivate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})

	component.child.Handle(tui.MsgActivate{Recipient: tui.Everyone})
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
}
