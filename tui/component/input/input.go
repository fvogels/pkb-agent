package input

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
)

type Component struct {
	size     tui.Size
	contents data.Value[string]
	style    tui.Style
	onChange func(string)
	label    *label.Component
}

func New(contents data.Value[string], style tui.Style, onChange func(string)) *Component {
	return &Component{
		contents: contents,
		style:    style,
		onChange: onChange,
		label:    label.New(contents, style),
	}
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
		component.onChange(updatedContents)
	} else if message.Key == "Backspace" {
		currentContents := component.contents.Get()

		if len(currentContents) > 0 {
			component.onChange(currentContents[:len(currentContents)-1])
		}
	}
}
