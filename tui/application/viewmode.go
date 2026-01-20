package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/border"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/keyview"
	"pkb-agent/tui/component/linksview"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
	"pkb-agent/tui/model"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type viewMode struct {
	tui.ComponentBase
	application                 *Application
	statusBar                   tui.Component
	highlightedNode             data.Value[*pkg.Node]
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newViewMode(application *Application) *viewMode {
	messageQueue := application.messageQueue

	selectedNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.SelectedNodes })
	intersectionNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.IntersectionNodes })
	highlightedNodeIndex := data.MapValue(&application.model, func(m *model.Model) int { return m.HighlightedNodeIndex })
	linkViewActive := data.MapValue(&application.model, func(m *model.Model) bool { return m.ShowNodeLinks })

	nodesView := nodeselection.New(messageQueue, selectedNodes, intersectionNodes, highlightedNodeIndex)
	statusBar := keyview.New(messageQueue, "status bar", application.bindings.all)
	highlightedNode := data.MapValue3(
		highlightedNodeIndex,
		intersectionNodes,
		selectedNodes,
		func(highlightedNodeIndex int, intersectionNodes list.List[*pkg.Node], selectedNodes list.List[*pkg.Node]) *pkg.Node {
			if intersectionNodes.Size() > 0 {
				return intersectionNodes.At(highlightedNodeIndex)
			} else if selectedNodes.Size() > 0 {
				return selectedNodes.At(selectedNodes.Size() - 1)
			} else {
				return nil
			}
		},
	)
	highlightedNodeViewer := data.MapValue2(
		highlightedNode,
		linkViewActive,
		func(highlightedNode *pkg.Node, linkViewActive bool) tui.Component {
			if highlightedNode != nil {
				if linkViewActive {
					return linksview.New(messageQueue, highlightedNode)
				} else {
					return highlightedNode.GetViewer(messageQueue)
				}
			} else {
				return nil
			}
		},
	)
	highlightedNodeViewerHolder := holder.New(messageQueue, highlightedNodeViewer)

	borderStyle := tcell.StyleDefault.Foreground(color.Reset).Background(color.Reset)
	root := docksouth.New(
		messageQueue,
		"view:docksouth[main|statusbar]",
		docknorth.New(
			messageQueue,
			"view:docknorth[nodes|nodeviewer]",
			border.New(messageQueue, nodesView, &borderStyle),
			border.New(messageQueue, highlightedNodeViewerHolder, &borderStyle),
			20,
		),
		statusBar,
		1,
	)

	nodesView.SetOnSelectionChanged(func(value int) {
		application.highlight(value)
	})

	result := viewMode{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "view mode",
			MessageQueue: messageQueue,
		},
		application: application,
		statusBar:   statusBar,
		root:        root,
	}

	return &result
}

func (component *viewMode) Render() tui.Grid {
	return component.root.Render()
}

func (component *viewMode) Handle(message tui.Message) {
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

func (component *viewMode) onStateUpdated() {
	messageQueue := component.application.messageQueue
	messageQueue.Enqueue(messages.MsgSetModeKeyBindings{
		Bindings: list.FromItems(
			BindingQuit,
			BindingSelect,
			BindingUnselect,
			BindingSearch,
			BindingSwitchLinksView,
			BindingLockNodes,
			BindingUnlockNodes,
		),
	})
}

func (component *viewMode) onKey(message tui.MsgKey) {
	application := component.application
	activeBindings := []tui.KeyBinding{
		BindingQuit,
		BindingSelect,
		BindingUnselect,
		BindingSearch,
		BindingSwitchLinksView,
		BindingLockNodes,
		BindingUnlockNodes,
	}

	if !tui.HandleKeyBindings(application.messageQueue, message, activeBindings...) {
		component.root.Handle(message)
	}
}
