package snippetview

import (
	"pkb-agent/tui"
	"pkb-agent/tui/component/ansiview"
	"pkb-agent/tui/data"
	"pkb-agent/util/syntaxhighlighting"
)

type Component struct {
	size            tui.Size
	source          data.Value[Source]
	formattedSource data.Variable[string]
	ansiView        *ansiview.Component
}

type Source struct {
	Contents string
	Language string
}

func New(messageQueue tui.MessageQueue, source data.Value[Source]) *Component {
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
	source := component.source.Get().Contents
	language := component.source.Get().Language

	reformatted, err := syntaxhighlighting.Highlight(source, language)
	if err != nil {
		panic("failed to render snippet")
	}

	component.formattedSource.Set(reformatted)
}
