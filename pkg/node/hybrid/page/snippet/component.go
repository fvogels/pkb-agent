package snippetpage

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/snippetview"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type pageComponent struct {
	tui.ComponentBase
	parent *Page
	child  *snippetview.Component
}

func NewPageComponent(messageQueue tui.MessageQueue, parent *Page) *pageComponent {
	source := data.NewConstant(snippetview.Source{
		Contents: parent.source,
		Language: parent.language,
	})

	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless snippet page",
			MessageQueue: messageQueue,
		},
		parent: parent,
		child:  snippetview.New(messageQueue, source),
	}

	return &component
}

func (component *pageComponent) Render() grid.FiniteGrid {
	return component.child.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgStateUpdated:
		component.child.Handle(message)

	default:
		component.child.Handle(message)
	}
}
