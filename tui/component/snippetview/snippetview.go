package snippetview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/ansiview"
	"pkb-agent/tui/data"
	"pkb-agent/util/syntaxhighlighting"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	source          data.Value[Source]
	formattedSource data.Variable[string]
	child           *ansiview.Component
}

type Source struct {
	Contents string
	Language string
}

func New(messageQueue tui.MessageQueue, source data.Value[Source]) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed snippet viewer",
			MessageQueue: messageQueue,
		},
		source:          source,
		formattedSource: data.NewVariable(""),
	}
	component.child = ansiview.New(messageQueue, &component.formattedSource)

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgStateUpdated:
		component.onStateUpdated()
		component.child.Handle(message)

	default:
		component.child.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
	component.reformatMarkdown()
}

func (component *Component) onStateUpdated() {
	component.reformatMarkdown()
}

func (component *Component) reformatMarkdown() {
	source := component.source.Get().Contents
	language := component.source.Get().Language

	reformatted, err := syntaxhighlighting.Highlight(source, language)
	if err != nil {
		panic("failed to render snippet")
	}

	component.formattedSource.Set(reformatted)
}
