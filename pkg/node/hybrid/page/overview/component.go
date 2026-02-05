package overview

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/tui"
	"pkb-agent/tui/component/numberedstringlist"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util"
	"pkb-agent/util/uid"
)

type pageComponent struct {
	tui.ComponentBase
	parent *Page
	root   tui.Component
}

func NewPageComponent(messageQueue tui.MessageQueue, parent *Page) *pageComponent {
	component := pageComponent{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless overview page",
			MessageQueue: messageQueue,
		},
		root: createRoot(messageQueue, parent.pages),
	}

	return &component
}

func createRoot(messageQueue tui.MessageQueue, pages []page.Page) tui.Component {
	captionsAsSlice := util.Map(pages, func(page page.Page) string {
		return page.GetCaption()
	})
	captionsAsList := list.FromSlice(captionsAsSlice)
	return numberedstringlist.New(messageQueue, data.NewConstant(captionsAsList))
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
