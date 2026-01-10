package input

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
)

type Component struct {
	size     tui.Size
	contents data.Value[string]
	style    *tui.Style
	onChange func(string)
	label    *label.Component
}

func New(contents data.Value[string]) *Component {
	style := tcell.StyleDefault
	subComponent := label.New(contents)
	subComponent.SetStyle(&style)

	return &Component{
		contents: contents,
		style:    &style,
		onChange: nil,
		label:    subComponent,
	}
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
	}
}

func (component *Component) Render() tui.Grid {
	return component.label.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size

	component.label.Handle(message)
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
