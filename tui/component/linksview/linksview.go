package linksview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	root tui.Component
	node *pkg.Node
}

func New(messageQueue tui.MessageQueue, node *pkg.Node) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed links view",
			MessageQueue: messageQueue,
		},
		root: label.New(messageQueue, "links view", data.NewConstant("links view")),
		node: node,
	}

	return &component
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgStateUpdated:
		component.root.Handle(message)
		component.onStateUpdated()

	default:
		component.root.Handle(message)
	}
}

func (component *Component) onStateUpdated() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.FromItems[tui.KeyBinding](
			tui.KeyBinding{
				Key:         "@",
				Description: "tralala",
				Message:     nil,
			},
		),
	})
}
