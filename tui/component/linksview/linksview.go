package linksview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	root tui.Component
}

func New(messageQueue tui.MessageQueue) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed links view",
			MessageQueue: messageQueue,
		},
		root: label.New(messageQueue, "links view", data.NewConstant("links view")),
	}

	return &component
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) Handle(message tui.Message) {
	component.root.Handle(message)
}
