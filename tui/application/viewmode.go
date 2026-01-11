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
	application *Application
	statusBar   *label.Component
	nodes       *nodeselection.Component
	root        tui.Component
}

func newViewMode(application *Application) *viewMode {
	model := &application.model

	nodesView := nodeselection.New(model.selectedNodes, model.intersectionNodes, model.highlightedNodeIndex)
	activeNodeViewer := holder.New(application.model.highlightedNodeViewer)
	caption := data.NewConstant("hello")
	statusBar := label.New("view:statusbar", caption)
	root := docksouth.New(
		docknorth.New(
			"view:docknorth[nodes|nodeviewer]",
			nodesView,
			activeNodeViewer,
			20,
		),
		statusBar,
		1,
	)

	nodesView.SetOnSelectionChanged(func(value int) { model.highlightedNodeIndex.Set(value) })

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
