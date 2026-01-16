package empty

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
)

type pageComponent struct {
	tui.ComponentBase
	child *label.Component
}

type msgCopySnippet struct{}

func NewPageComponent(messageQueue tui.MessageQueue) *pageComponent {
	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless snippet page",
			MessageQueue: messageQueue,
		},
		child: label.New(messageQueue, "empty page label", data.NewConstant("no pages")),
	}

	return &component
}

func (component *pageComponent) Render() tui.Grid {
	return component.child.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		component.onKey(message)

	case tui.MsgActivate:
		if message.ShouldRespond(component.Identifier) {
			component.onActivate()
		}

	default:
		component.child.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	// No key bindings for this page
}

func (component *pageComponent) onActivate() {
	component.MessageQueue.Enqueue(page.MsgSetPageKeyBindings{
		Bindings: list.FromItems[tui.KeyBinding](),
	})
}
