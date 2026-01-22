package input

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	contents data.Value[string]
	style    *tui.Style
	onChange func(string)
	child    *label.Component
}

func New(messageQueue tui.MessageQueue, contents data.Value[string]) *Component {
	style := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)
	child := label.New(messageQueue, "input[label]", contents)
	child.SetStyle(&style)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "nameless label",
			MessageQueue: messageQueue,
		},
		contents: contents,
		style:    &style,
		onChange: nil,
		child:    child,
	}

	return &component
}

func (component *Component) SetStyle(style *tui.Style) {
	component.style = style
}

func (component *Component) SetOnChange(callback func(string)) {
	component.onChange = callback
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgKey:
		component.onKey(message)

	default:
		component.child.Handle(message)
	}
}

func (component *Component) Render() grid.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	component.child.Handle(message)
}

func (component *Component) onKey(message tui.MsgKey) {
	if len(message.Key) == 1 {
		updatedContents := component.contents.Get() + message.Key
		component.signalOnChange(updatedContents)
	} else if message.Key == "Backspace" {
		currentContents := component.contents.Get()

		if len(currentContents) > 0 {
			component.signalOnChange(currentContents[:len(currentContents)-1])
		}
	}
}

func (component *Component) signalOnChange(contents string) {
	if component.onChange != nil {
		component.onChange(contents)
	}
}
