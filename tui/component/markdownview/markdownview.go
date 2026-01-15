package markdownview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/ansiview"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
	"pkb-agent/util/markdown"
)

type Component struct {
	tui.ComponentBase
	source          data.Value[string]
	formattedSource data.Variable[string]
	ansiView        *ansiview.Component
}

func New(messageQueue tui.MessageQueue, source data.Value[string]) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed markdown viewer",
			MessageQueue: messageQueue,
		},
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

	case tui.MsgActivate:
		if message.ShouldRespond(component.Identifier) {
			component.onActivate()
		}
	}
}

func (component *Component) Render() tui.Grid {
	return component.ansiView.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.reformatMarkdown()
}

func (component *Component) reformatMarkdown() {
	reformatted, err := markdown.Render(component.source.Get(), component.Size.Width)
	if err != nil {
		panic("failed to render markdown")
	}

	component.formattedSource.Set(reformatted)
}

func (component *Component) onActivate() {
	component.ansiView.Handle(tui.MsgActivate{Recipient: tui.Everyone})
}
