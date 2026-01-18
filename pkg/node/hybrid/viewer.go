package hybrid

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/pkg/node/hybrid/page/empty"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	rawNode                *RawNode
	data                   *nodeData // (strong) pointer to the node data, keeps information alive while viewer exists
	activePageIndex        data.Variable[int]
	pageViewers            []tui.Component
	actionKeyBindings      data.Variable[list.List[tui.KeyBinding]] // Key bindings associated with the node
	pageKeyBindings        data.Variable[list.List[tui.KeyBinding]] // Page specific key bindings
	keyBindings            data.Value[list.List[tui.KeyBinding]]    // Concatenation of action key bindings and page key bindings
	activePageViewer       data.Value[tui.Component]
	activePageViewerHolder holder.Component
	pageStatus             data.Value[string]
	pageStatusView         tui.Component
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
	component.pageViewers = component.createPageViewers(messageQueue, pages)

	component.actionKeyBindings = data.NewVariable(list.New[tui.KeyBinding]())
	component.pageKeyBindings = data.NewVariable(list.New[tui.KeyBinding]())

	// keyBindings should be kept equal to the concatenation of actionKeyBindings and pageKeyBindings
	component.keyBindings = data.MapValue2(
		&component.actionKeyBindings,
		&component.pageKeyBindings,
		func(xs, ys list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			return list.Concatenate(xs, ys)
		},
	)

	// Whenever keyBindings change, send a message
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
	component.pageStatus = data.MapValue(
		&component.activePageIndex,
		func(pageIndex int) string {
			if len(pages) > 0 {
				return fmt.Sprintf("Page %d/%d: %s", pageIndex+1, len(pages), pages[pageIndex].GetCaption())
			} else {
				return "No pages"
			}
		},
	)
	component.pageStatusView = label.New(
		messageQueue,
		"page status",
		component.pageStatus,
	)

	component.root = docksouth.New(
		messageQueue,
		"docksouth[page|pagestatus]",
		&component.activePageViewerHolder,
		component.pageStatusView,
		1,
	)

	return &component
}

func (component *Component) createPageViewers(messageQueue tui.MessageQueue, pages []page.Page) []tui.Component {
	pageViewers := make([]tui.Component, len(pages))

	for pageIndex, page := range pages {
		viewer := page.CreateViewer(messageQueue)
		pageViewers[pageIndex] = viewer
	}

	return pageViewers
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
