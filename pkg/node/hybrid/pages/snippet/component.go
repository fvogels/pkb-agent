package snippetpage

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/snippetview"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"

	"golang.design/x/clipboard"
)

type pageComponent struct {
	tui.ComponentBase
	parent        *Page
	snippetViewer *snippetview.Component
	bindingCopy   tui.KeyBinding
}

type msgCopySnippet struct{}

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
		parent:        parent,
		snippetViewer: snippetview.New(messageQueue, source),
		bindingCopy: tui.KeyBinding{
			Key:         "c",
			Description: "copy",
			Message:     msgCopySnippet{},
		},
	}

	return &component
}

func (component *pageComponent) Render() tui.Grid {
	return component.snippetViewer.Render()
}

func (component *pageComponent) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		component.onKey(message)

	case msgCopySnippet:
		component.onCopySnippet()

	default:
		component.snippetViewer.Handle(message)
	}
}

func (component *pageComponent) onKey(message tui.MsgKey) {
	tui.HandleKeyBindings(component.MessageQueue, message, component.bindingCopy)
}

func (component *pageComponent) onCopySnippet() {
	clipboard.Write(clipboard.FmtText, ([]byte)(component.parent.source))
}
