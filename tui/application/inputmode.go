package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/pkg"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/docknorth"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/holder"
	"pkb-agent/tui/component/input"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
	"pkb-agent/tui/model"
	"pkb-agent/ui/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type inputMode struct {
	tui.ComponentBase
	application                 *Application
	inputField                  *input.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newInputMode(application *Application) *inputMode {
	messageQueue := application.messageQueue

	selectedNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.SelectedNodes })
	intersectionNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.IntersectionNodes })
	highlightedNodeIndex := data.MapValue(&application.model, func(m *model.Model) int { return m.HighlightedNodeIndex })
	searchString := data.MapValue(&application.model, func(m *model.Model) string { return m.Input })

	nodesView := nodeselection.New(
		messageQueue,
		selectedNodes,
		intersectionNodes,
		highlightedNodeIndex,
	)
	nodesView.SetOnSelectionChanged(func(value int) {
		application.highlight(value)
	})

	highlightedNodeViewer := data.MapValue3(
		highlightedNodeIndex,
		intersectionNodes,
		selectedNodes,
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

	inputField := input.New(messageQueue, searchString)
	style := tcell.StyleDefault.Background(color.Red)
	inputField.SetStyle(&style)
	inputField.SetOnChange(func(s string) { application.updateInputAndHighlightBestMatch(s) })

	root := docksouth.New(
		messageQueue,
		"input:[main|input]",
		docknorth.New(
			messageQueue,
			"input:docknorth[nodes|nodeviewer]",
			nodesView,
			highlightedNodeViewerHolder,
			30,
		),
		inputField,
		1,
	)

	result := inputMode{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "input mode",
			MessageQueue: messageQueue,
		},
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

	case tui.MsgActivate:
		if message.ShouldRespond(mode.Identifier) {
			mode.onActivate()
		}

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

func (mode *inputMode) onActivate() {
	mode.application.messageQueue.Enqueue(messages.MsgSetModeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})

	mode.root.Handle(tui.MsgActivate{Recipient: tui.Everyone})
}
