package hybrid

import (
	"fmt"
	"log/slog"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node"
	"pkb-agent/pkg/node/hybrid/page"
	"pkb-agent/pkg/node/hybrid/page/empty"
	"pkb-agent/pkg/node/hybrid/page/overview"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	rawNode                *RawNode
	data                   *nodeData // (strong) pointer to the node data, keeps information alive while viewer exists
	activePageIndex        data.Variable[int]
	pages                  []page.Page
	pageViewers            []tui.Component
	bindings               keyBindings
	activePageViewer       data.Value[tui.Component]
	activePageViewerHolder holder.Component
	pageStatus             data.Value[string]
	pageStatusView         tui.Component
	root                   tui.Component
}

type keyBindings struct {
	actions    []tui.KeyBinding                         // Key bindings associated with the node
	jumpToPage []tui.KeyBinding                         // Key bindings to jump to page
	page       data.Variable[list.List[tui.KeyBinding]] // Page specific key bindings
	all        data.Value[list.List[tui.KeyBinding]]    // Concatenation of action key bindings and page key bindings
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

	component.pages = addOverviewPage(nodeData.pages)
	component.pageViewers = component.createPageViewers(messageQueue, component.pages)

	component.createKeyBindings(nodeData, &component.bindings)

	if len(component.pages) == 0 {
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
			if pageIndex == 0 {
				// First page is overview, do not mention index
				return component.pages[0].GetCaption()
			} else {
				totalPages := len(component.pages) - 1 // Subtract one for overview page
				caption := component.pages[pageIndex].GetCaption()

				return fmt.Sprintf(
					"Page %d/%d: %s",
					pageIndex,
					totalPages,
					caption,
				)
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

func addOverviewPage(pages []page.Page) []page.Page {
	overviewPage := overview.New(pages)

	return append([]page.Page{overviewPage}, pages...)
}

func (component *Component) createKeyBindings(nodeData *nodeData, bindings *keyBindings) {
	bindings.actions = component.createActionKeyBindings(nodeData.actions)
	bindings.jumpToPage = component.createJumpToPageKeyBindings()
	bindings.page = data.NewVariable(list.New[tui.KeyBinding]())

	bindings.all = data.MapValue(
		&bindings.page,
		func(pageBindings list.List[tui.KeyBinding]) list.List[tui.KeyBinding] {
			return list.Concatenate(
				list.FromSlice(bindings.actions),
				list.FromSlice(bindings.jumpToPage),
				pageBindings,
			)
		},
	)
}

func (component *Component) createJumpToPageKeyBindings() []tui.KeyBinding {
	keys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	pageCount := len(component.pages)
	keyBindings := make([]tui.KeyBinding, 0, pageCount)

	for pageIndex := range pageCount {
		var description string

		if pageIndex == 0 {
			description = "Overview"
		} else {
			description = fmt.Sprintf("To page %d", pageIndex)
		}

		keyBinding := tui.KeyBinding{
			Key:         keys[pageIndex],
			Description: description,
			Message:     page.MsgSetActivePage{PageIndex: pageIndex},
		}

		keyBindings = append(keyBindings, keyBinding)
	}

	return keyBindings
}

func (component *Component) signalNodeKeyBindingsUpdate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: component.bindings.all.Get(),
	})
}

func (component *Component) createActionKeyBindings(actions []node.Action) []tui.KeyBinding {
	keys := []rune{'Q', 'W', 'E', 'R', 'T', 'Y', 'U', 'I', 'O', 'P'}
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

	case page.MsgSetActivePage:
		component.onSetActivePage(message)

	default:
		component.root.Handle(message)
	}
}

func (component *Component) onSetPageKeyBindings(message page.MsgSetPageKeyBindings) {
	component.bindings.page.Set(message.Bindings)
	component.signalNodeKeyBindingsUpdate()
}

func (component *Component) onSetActivePage(message page.MsgSetActivePage) {
	component.withActivePage(func(page page.Page, viewer tui.Component) {
		component.setActivePage(message.PageIndex)
	})
	component.Handle(tui.MsgStateUpdated{})
	component.Handle(tui.MsgResize{Size: component.Size})
}

func (component *Component) Render() grid.FiniteGrid {
	slog.Debug("Rendering hybrid node viewer", slog.String("size", component.Size.String()))

	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	component.root.Handle(message)
}

func (component *Component) withActivePage(f func(page page.Page, viewer tui.Component)) {
	if len(component.pageViewers) > 0 {
		activePage := component.pages[component.activePageIndex.Get()]
		activeViewer := component.pageViewers[component.activePageIndex.Get()]

		f(activePage, activeViewer)
	}
}

func (component *Component) onKey(message tui.MsgKey) {
	if tui.HandleKeyBindings(component.MessageQueue, message, component.bindings.actions...) {
		return
	}

	if tui.HandleKeyBindings(component.MessageQueue, message, component.bindings.jumpToPage...) {
		return
	}

	switch message.Key {
	case "Tab":
		newActivePageIndex := (component.activePageIndex.Get() + 1) % len(component.pageViewers)
		component.Handle(page.MsgSetActivePage{PageIndex: newActivePageIndex})

	case "Backtab":
		newActivePageIndex := (component.activePageIndex.Get() - 1 + len(component.pageViewers)) % len(component.pageViewers)
		component.Handle(page.MsgSetActivePage{PageIndex: newActivePageIndex})

	default:
		component.withActivePage(func(page page.Page, viewer tui.Component) {
			viewer.Handle(message)
		})
	}
}

func (component *Component) setActivePage(index int) {
	component.activePageIndex.Set(index)
}
