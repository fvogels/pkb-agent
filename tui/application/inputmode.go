package application

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/docksouth"
	"pkb-agent/tui/component/input"
	"pkb-agent/tui/component/nodeselection"
	"strings"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type inputMode struct {
	application *Application
	inputField  *input.Component
	nodes       *nodeselection.Component
	root        tui.Component
}

func newInputMode(application *Application) *inputMode {
	model := &application.model

	nodesView := nodeselection.New(model.selectedNodes, model.intersectionNodes, model.highlightedNodeIndex)
	nodesView.SetOnSelectionChanged(func(value int) { model.highlightedNodeIndex.Set(value) })

	inputField := input.New(model.input)
	style := tcell.StyleDefault.Background(color.Red)
	inputField.SetStyle(&style)
	inputField.SetOnChange(func(s string) { model.input.Set(strings.ToLower(s)) })

	root := docksouth.New(nodesView, inputField, 1)

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

	default:
		mode.root.Handle(message)
	}
}

func (mode *inputMode) onKey(message tui.MsgKey) {
	switch message.Key {
	case "Enter":
		mode.application.selectHighlightedNode()
		mode.application.model.input.Set("")
		mode.application.switchMode(mode.application.viewMode)

	default:
		mode.root.Handle(message)
	}
}
