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
	"pkb-agent/tui/component/input"
	"pkb-agent/tui/component/nodeselection"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/tui/model"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type searchMode struct {
	tui.ComponentBase
	application                 *Application
	inputField                  *input.Component
	highlightedNodeViewer       data.Value[tui.Component]
	highlightedNodeViewerHolder *holder.Component
	nodes                       *nodeselection.Component
	root                        tui.Component
}

func newSearchMode(application *Application) *searchMode {
	messageQueue := application.messageQueue

	selectedNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.SelectedNodes })
	intersectionNodes := data.MapValue(&application.model, func(m *model.Model) list.List[*pkg.Node] { return m.IntersectionNodes })
	highlightedNodeIndex := data.MapValue(&application.model, func(m *model.Model) int { return m.HighlightedNodeIndex })
	searchString := data.MapValue(&application.model, func(m *model.Model) string { return m.Input })
	lockCount := data.MapValue(&application.model, func(m *model.Model) int { return m.LockedNodeCount })

	nodesView := nodeselection.New(
		messageQueue,
		selectedNodes,
		intersectionNodes,
		highlightedNodeIndex,
		lockCount,
	)
	nodesView.SetOnSelectionChanged(func(value int) {
		application.highlight(value)
	})

	highlightedNodeViewer := data.Cache(
		data.MapValue3(
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
		),
	)
	highlightedNodeViewerHolder := holder.New(messageQueue, highlightedNodeViewer)

	inputField := input.New(messageQueue, searchString)
	style := tcell.StyleDefault.Background(color.Red)
	inputField.SetStyle(&style)
	inputField.SetOnChange(func(s string) { application.updateInputAndHighlightBestMatch(s) })

	borderStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	root := docksouth.New(
		messageQueue,
		"input:[main|input]",
		docknorth.New(
			messageQueue,
			"input:docknorth[nodes|nodeviewer]",
			border.New(messageQueue, nodesView, &borderStyle),
			border.New(messageQueue, highlightedNodeViewerHolder, &borderStyle),
			20,
		),
		inputField,
		1,
	)

	result := searchMode{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "search mode",
			MessageQueue: messageQueue,
		},
		application: application,
		inputField:  inputField,
		root:        root,
		nodes:       nodesView,
	}

	return &result
}

func (component *searchMode) Render() grid.Grid {
	return component.root.Render()
}

func (component *searchMode) Handle(message tui.Message) {
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

func (component *searchMode) onKey(message tui.MsgKey) {
	switch message.Key {
	case "Enter":
		component.application.selectHighlightedAndClearInput()
		component.application.switchMode(component.application.mode.view)

	case "Esc":
		component.application.clearInput()
		component.application.switchMode(component.application.mode.view)

	default:
		component.inputField.Handle(message)
		component.nodes.Handle(message)
	}
}

func (component *searchMode) onStateUpdated() {
	component.application.messageQueue.Enqueue(messages.MsgSetModeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})
}
