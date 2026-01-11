package application

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
)

type viewMode struct {
	application                 *Application
	statusBar                   *label.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newViewMode(application *Application) *viewMode {
	model := &application.model

	nodesView := nodeselection.New(model.selectedNodes, model.intersectionNodes, model.highlightedNodeIndex)
	caption := data.NewConstant("hello")
	statusBar := label.New("view:statusbar", caption)
	highlightedNodeViewer := data.NewVariable[tui.Component](nil)
	highlightedNodeViewerHolder := holder.New(highlightedNodeViewer)

	data.DefineReaction(func() {
		var viewer tui.Component

		if model.intersectionNodes.Size() > 0 {
			viewer = model.intersectionNodes.At(model.highlightedNodeIndex.Get()).GetViewer()
		} else if model.selectedNodes.Size() > 0 {
			viewer = model.selectedNodes.At(model.selectedNodes.Size() - 1).GetViewer()
		} else {
			// Should not happen
			viewer = nil
		}

		highlightedNodeViewer.Set(viewer)

	}, model.highlightedNodeIndex, model.selectedNodes)

	root := docksouth.New(
		"view:docksouth[main|statusbar]",
		docknorth.New(
			"view:docknorth[nodes|nodeviewer]",
			nodesView,
			highlightedNodeViewerHolder,
			20,
		),
		statusBar,
		1,
	)

	nodesView.SetOnSelectionChanged(func(value int) { model.highlightedNodeIndex.Set(value) })

	// // Cause node viewer to be updated automatically
	// updateHighlightedNodeViewer := func() {
	// 	var viewer tui.Component

	// 	if intersectionNodes.Size() > 0 {
	// 		viewer = intersectionNodes.At(highlightedNodeIndex.Get()).GetViewer()
	// 	} else if selectedNodes.Size() > 0 {
	// 		viewer = selectedNodes.At(selectedNodes.Size() - 1).GetViewer()
	// 	} else {
	// 		// Should not happen
	// 		viewer = nil
	// 	}

	// 	highlightedNodeViewer.Set(viewer)
	// }
	// updateHighlightedNodeViewer()
	// data.DefineReaction(updateHighlightedNodeViewer, highlightedNodeIndex, selectedNodes)

	result := viewMode{
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

	default:
		mode.root.Handle(message)
	}
}

func (mode *viewMode) onKey(message tui.MsgKey) {
	application := mode.application

	switch message.Key {
	case "q":
		application.running = false

	case "Enter":
		application.selectHighlightedNode()

	case "Delete":
		application.unselectLastNode()

	case "/":
		application.switchMode(mode.application.inputMode)

	default:
		mode.root.Handle(message)
	}
}

// func (mode *viewMode) updateLayout() {
// 	message := tui.MsgResize{
// 		Size: mode.application.size,
// 	}
// 	mode.root.Handle(message)
// }
