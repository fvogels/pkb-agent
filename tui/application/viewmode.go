package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/keyview"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
)

type viewMode struct {
	tui.ComponentBase
	application                 *Application
	statusBar                   tui.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newViewMode(application *Application) *viewMode {
	model := &application.model
	messageQueue := application.messageQueue

	nodesView := nodeselection.New(messageQueue, model.SelectedNodes(), model.IntersectionNodes(), model.HighlightedNodeIndex())
	statusBar := keyview.New(messageQueue, "status bar", application.keyBindings)
	highlightedNodeViewer := data.MapValue3(
		model.HighlightedNodeIndex(),
		model.IntersectionNodes(),
		model.SelectedNodes(),
		func(highlightedNodeIndex int, intersectionNodes list.List[*pkg.Node], selectedNodes list.List[*pkg.Node]) tui.Component {
			if intersectionNodes.Size() > 0 {
				return intersectionNodes.At(highlightedNodeIndex).GetViewer(messageQueue)
			} else if selectedNodes.Size() > 0 {
				return selectedNodes.At(selectedNodes.Size() - 1).GetViewer(messageQueue)
			} else {
				return nil
			}
		},
	)
	highlightedNodeViewerHolder := holder.New(messageQueue, highlightedNodeViewer)

	root := docksouth.New(
		messageQueue,
		"view:docksouth[main|statusbar]",
		docknorth.New(
			messageQueue,
			"view:docknorth[nodes|nodeviewer]",
			nodesView,
			highlightedNodeViewerHolder,
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

func (mode *viewMode) Render() tui.Grid {
	return mode.root.Render()
}

func (mode *viewMode) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		mode.onKey(message)

	case tui.MsgActivate:
		if message.ShouldRespond(mode.Identifier) {
			mode.onActivate()
		}

	default:
		mode.root.Handle(message)
	}
}

func (mode *viewMode) onActivate() {
	messageQueue := mode.application.messageQueue
	messageQueue.Enqueue(messages.MsgSetModeKeyBindings{
		Bindings: list.FromItems(
			BindingQuit,
			BindingSelect,
			BindingUnselect,
			BindingSearch,
		),
	})

	mode.root.Handle(tui.MsgActivate{Recipient: tui.Everyone})
}

func (mode *viewMode) onKey(message tui.MsgKey) {
	application := mode.application
	activeBindings := []tui.KeyBinding{
		BindingQuit,
		BindingSelect,
		BindingUnselect,
		BindingSearch,
	}

	if !tui.HandleKeyBindings(application.messageQueue, message, activeBindings...) {
		mode.root.Handle(message)
	}
}
