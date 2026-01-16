package hybrid

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/pkg/node/hybrid/page/empty"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	rawNode                *RawNode
	data                   *nodeData // (strong) pointer to the node data, keeps information alive while viewer exists
	activePageIndex        data.Variable[int]
	pageViewers            []tui.Component
	actionKeyBindings      data.Variable[list.List[tui.KeyBinding]]
	pageKeyBindings        data.Variable[list.List[tui.KeyBinding]]
	keyBindings            data.Value[list.List[tui.KeyBinding]]
	activePageViewer       data.Value[tui.Component]
	activePageViewerHolder holder.Component
	root                   tui.Component
}

func NewViewer(messageQueue tui.MessageQueue, rawNode *RawNode, nodeData *nodeData) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed hybrid node viewer",
			MessageQueue: messageQueue,
		},
		rawNode:         rawNode,
		activePageIndex: data.NewVariable(0),
		data:            nodeData,
	}

	pages := nodeData.pages
	component.pageViewers = make([]tui.Component, len(pages))

	for pageIndex, page := range pages {
		viewer := page.CreateViewer(messageQueue)
		component.pageViewers[pageIndex] = viewer
	}

	component.actionKeyBindings = data.NewVariable(list.New[tui.KeyBinding]())
	component.pageKeyBindings = data.NewVariable(list.New[tui.KeyBinding]())
	component.keyBindings = data.MapValue2(
		&component.actionKeyBindings,
		&component.pageKeyBindings,
		func(xs, ys list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			return list.Concatenate(xs, ys)
		},
	)
	component.keyBindings.Observe(func() {
		messageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
			Bindings: component.keyBindings.Get(),
		})
	})

	if len(pages) == 0 {
		component.activePageViewer = data.NewConstant[tui.Component](empty.NewPageComponent(messageQueue))
	} else {
		component.activePageViewer = data.MapValue(&component.activePageIndex, func(index int) tui.Component {
			return component.pageViewers[index]
		})
	}
	component.activePageViewerHolder = *holder.New(messageQueue, component.activePageViewer)
	component.root = &component.activePageViewerHolder

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgActivate:
		message.Respond(
			component.Identifier,
			component.onActivate,
			component.root,
		)

	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgKey:
		component.onKey(message)

	case page.MsgSetPageKeyBindings:
		component.onSetPageKeyBindings(message)

	default:
		component.root.Handle(message)
	}
}

func (component *Component) onActivate() {
	if len(component.pageViewers) > 0 {
		component.setActivePage(0)
	}
}

func (component *Component) onSetPageKeyBindings(message page.MsgSetPageKeyBindings) {
	component.pageKeyBindings.Set(message.Bindings)
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	component.root.Handle(message)
}

func (component *Component) withActivePage(f func(page page.Page, viewer tui.Component)) {
	if len(component.pageViewers) > 0 {
		activePage := component.data.pages[component.activePageIndex.Get()]
		activeViewer := component.pageViewers[component.activePageIndex.Get()]

		f(activePage, activeViewer)
	}
}

func (component *Component) onKey(message tui.MsgKey) {
	switch message.Key {
	case "Tab":
		component.withActivePage(func(page page.Page, viewer tui.Component) {
			component.setActivePage((component.activePageIndex.Get() + 1) % len(component.pageViewers))
		})

	default:
		component.withActivePage(func(page page.Page, viewer tui.Component) {
			viewer.Handle(message)
		})
	}
}

func (component *Component) setActivePage(index int) {
	component.activePageIndex.Set(index)
}
