package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/keyview"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
)

type viewMode struct {
	application                 *Application
	statusBar                   tui.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newViewMode(application *Application) *viewMode {
	model := &application.model

	nodesView := nodeselection.New(model.SelectedNodes(), model.IntersectionNodes(), model.HighlightedNodeIndex())
	statusBar := keyview.New(application.messageQueue, "status bar", application.keyBindings)
	highlightedNodeViewer := data.MapValue3(
		model.HighlightedNodeIndex(),
		model.IntersectionNodes(),
		model.SelectedNodes(),
		func(highlightedNodeIndex int, intersectionNodes list.List[*pkg.Node], selectedNodes list.List[*pkg.Node]) tui.Component {
			if intersectionNodes.Size() > 0 {
				return intersectionNodes.At(highlightedNodeIndex).GetViewer(application.messageQueue)
			} else if selectedNodes.Size() > 0 {
				return selectedNodes.At(selectedNodes.Size() - 1).GetViewer(application.messageQueue)
			} else {
				return nil
			}
		},
	)
	highlightedNodeViewerHolder := holder.New(application.messageQueue, highlightedNodeViewer)

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

	nodesView.SetOnSelectionChanged(func(value int) {
		application.highlight(value)
	})

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
