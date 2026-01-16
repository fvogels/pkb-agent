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
	child           *ansiview.Component
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
	component.child = ansiview.New(messageQueue, &component.formattedSource)
	component.source.Observe(func() { component.reformatMarkdown() })

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgActivate:
		message.Respond(
			component.Identifier,
			func() {},
			component.child,
		)
	}
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
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
