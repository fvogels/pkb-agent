package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/input"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type inputMode struct {
	application                 *Application
	inputField                  *input.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newInputMode(application *Application) *inputMode {
	model := &application.model

	nodesView := nodeselection.New(model.SelectedNodes(), model.IntersectionNodes(), model.HighlightedNodeIndex())
	nodesView.SetOnSelectionChanged(func(value int) {
		application.highlight(value)
	})

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

	inputField := input.New(model.Input())
	style := tcell.StyleDefault.Background(color.Red)
	inputField.SetStyle(&style)
	inputField.SetOnChange(func(s string) { application.updateInputAndHighlightBestMatch(s) })

	root := docksouth.New(
		"input:[main|input]",
		docknorth.New(
			"input:docknorth[nodes|nodeviewer]",
			nodesView,
			highlightedNodeViewerHolder,
			30,
		),
		inputField,
		1,
	)

	result := inputMode{
		application: application,
		inputField:  inputField,
		root:        root,
	}

	return &result
}

func (mode *inputMode) Render() tui.Grid {
	return mode.root.Render()
}

func (mode *inputMode) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgKey:
		mode.onKey(message)

	case MsgActivateMode:
		mode.onActivateMode()

	default:
		mode.root.Handle(message)
	}
}

func (mode *inputMode) onKey(message tui.MsgKey) {
	switch message.Key {
	case "Enter":
		mode.application.selectHighlightedAndClearInput()
		mode.application.switchMode(mode.application.viewMode)

	default:
		mode.root.Handle(message)
	}
}

func (mode *inputMode) onActivateMode() {
	mode.application.messageQueue.Enqueue(MsgSetModeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})
}
