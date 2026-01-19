package snippetpage

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/snippetview"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"

	"golang.design/x/clipboard"
)

type pageComponent struct {
	tui.ComponentBase
	parent      *Page
	child       *snippetview.Component
	bindingCopy tui.KeyBinding
}

type msgCopySnippet struct{}

func (message msgCopySnippet) String() string {
	return "msgCopySnippet"
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
		bindingCopy: tui.KeyBinding{
			Key:         "c",
			Description: "copy",
			Message:     msgCopySnippet{},
		},
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

	case msgCopySnippet:
		component.onCopySnippet()

	case tui.MsgStateUpdated:
		component.child.Handle(message)
		component.onStateUpdated()

	default:
		component.child.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	tui.HandleKeyBindings(component.MessageQueue, message, component.bindingCopy)
}

func (component *pageComponent) onCopySnippet() {
	clipboard.Write(clipboard.FmtText, ([]byte)(component.parent.source))
}

func (component *pageComponent) onStateUpdated() {
	component.MessageQueue.Enqueue(page.MsgSetPageKeyBindings{
		Bindings: list.FromItems(
			component.bindingCopy,
		),
	})
}
