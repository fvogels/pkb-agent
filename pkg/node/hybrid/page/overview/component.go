package overview

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type pageComponent struct {
	tui.ComponentBase
	page *Page
	root tui.Component
}

func NewPageComponent(messageQueue tui.MessageQueue, page *Page) *pageComponent {
	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless overview page",
			MessageQueue: messageQueue,
		},
		root: label.New(messageQueue, "label", data.NewConstant(fmt.Sprintf("%d pages", len(page.pages)))),
	}

	return &component
}

func (component *pageComponent) Render() grid.FiniteGrid {
	return component.root.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		component.onKey(message)

	case tui.MsgStateUpdated:
		component.root.Handle(message)
		component.onStateUpdated()

	default:
		component.root.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	// No key bindings for this page
}

func (component *pageComponent) onStateUpdated() {
	component.MessageQueue.Enqueue(page.MsgSetPageKeyBindings{
		Bindings: list.FromItems[tui.KeyBinding](),
	})
}
