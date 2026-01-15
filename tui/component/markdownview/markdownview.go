package markdownview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/ansiview"
	"pkb-agent/tui/data"
	"pkb-agent/util/markdown"
)

type Component struct {
	size            tui.Size
	source          data.Value[string]
	formattedSource data.Variable[string]
	ansiView        *ansiview.Component
}

func New(messageQueue tui.MessageQueue, source data.Value[string]) *Component {
	component := Component{
		source:          source,
		formattedSource: data.NewVariable(""),
	}
	component.ansiView = ansiview.New(messageQueue, &component.formattedSource)
	component.source.Observe(func() { component.reformatMarkdown() })

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return component.ansiView.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
	component.reformatMarkdown()
}

func (component *Component) reformatMarkdown() {
	reformatted, err := markdown.Render(component.source.Get(), component.size.Width)
	if err != nil {
		panic("failed to render markdown")
	}

	component.formattedSource.Set(reformatted)
}
