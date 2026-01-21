package hybrid

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node"
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
	bindings               keyBindings
	activePageViewer       data.Value[tui.Component]
	activePageViewerHolder holder.Component
	pageStatus             data.Value[string]
	pageStatusView         tui.Component
	root                   tui.Component
}

type keyBindings struct {
	actions []tui.KeyBinding                         // Key bindings associated with the node
	page    data.Variable[list.List[tui.KeyBinding]] // Page specific key bindings
	all     data.Value[list.List[tui.KeyBinding]]    // Concatenation of action key bindings and page key bindings
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

	component.createKeyBindings(nodeData, &component.bindings)

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

func (component *Component) createKeyBindings(nodeData *nodeData, bindings *keyBindings) {
	bindings.actions = component.createActionKeyBindings(nodeData.actions)
	bindings.page = data.NewVariable(list.New[tui.KeyBinding]())

	bindings.all = data.MapValue(
		&bindings.page,
		func(pageBindings list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			return list.Concatenate(list.FromSlice(bindings.actions), pageBindings)
		},
	)
}

func (component *Component) signalNodeKeyBindingsUpdate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: component.bindings.all.Get(),
	})
}

func (component *Component) createActionKeyBindings(actions []node.Action) []tui.KeyBinding {
	keys := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	keyBindings := make([]tui.KeyBinding, len(actions))

	for index, action := range actions {
		actionCopy := action
		description := action.GetDescription()
		key := string(keys[index])
		keyBindings[index] = tui.KeyBinding{
			Key:         key,
			Description: description,
			Message: tui.MsgCommand{
				Command: func() {
					go func() {
						actionCopy.Perform()
					}()
				},
			},
		}
	}

	return keyBindings
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
	case tui.MsgStateUpdated:
		component.root.Handle(message)

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

func (component *Component) onSetPageKeyBindings(message page.MsgSetPageKeyBindings) {
	component.bindings.page.Set(message.Bindings)
	component.signalNodeKeyBindingsUpdate()
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
	if tui.HandleKeyBindings(component.MessageQueue, message, component.bindings.actions...) {
		return
	}

	switch message.Key {
	case "Tab":
		component.withActivePage(func(page page.Page, viewer tui.Component) {
			component.setActivePage((component.activePageIndex.Get() + 1) % len(component.pageViewers))
		})
		component.Handle(tui.MsgStateUpdated{})
		component.Handle(tui.MsgResize{Size: component.Size})

	default:
		component.withActivePage(func(page page.Page, viewer tui.Component) {
			viewer.Handle(message)
		})
	}
}

func (component *Component) setActivePage(index int) {
	component.activePageIndex.Set(index)
}
