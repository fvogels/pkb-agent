package markdownpage

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/markdownview"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"

	"golang.design/x/clipboard"
)

type pageComponent struct {
	tui.ComponentBase
	parent         *Page
	markdownViewer *markdownview.Component
}

type msgCopySnippet struct{}

func NewPageComponent(messageQueue tui.MessageQueue, parent *Page) *pageComponent {
	source := data.NewConstant(parent.source)

	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless snippet page",
			MessageQueue: messageQueue,
		},
		parent:         parent,
		markdownViewer: markdownview.New(messageQueue, source),
	}

	return &component
}

func (component *pageComponent) Render() tui.Grid {
	return component.markdownViewer.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		component.onKey(message)

	case msgCopySnippet:
		component.onCopySnippet()

	case page.MsgActivatePage:
		component.MessageQueue.Enqueue(page.MsgSetPageKeyBindings{
			Bindings: list.FromItems[tui.KeyBinding](),
		})

	default:
		component.markdownViewer.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	// tui.HandleKeyBindings(component.MessageQueue, message, component.bindingCopy)
}

func (component *pageComponent) onCopySnippet() {
	clipboard.Write(clipboard.FmtText, ([]byte)(component.parent.source))
}
