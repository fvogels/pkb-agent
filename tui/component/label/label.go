package label

import (
	"log/slog"
	"pkb-agent/tui"
	"pkb-agent/tui/data"

	"github.com/gdamore/tcell/v3"
)

type Component struct {
	name     string
	size     tui.Size
	contents data.Value[string]
	style    *tui.Style
}

func New(name string, contents data.Value[string]) *Component {
	return &Component{
		name:     name,
		contents: contents,
		style:    &tcell.StyleDefault,
	}
}

func (component *Component) SetStyle(style *tui.Style) {
	component.style = style
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	return newGrid(component, []rune(component.contents.Get()))
}

func (component *Component) onResize(message tui.MsgResize) {
	slog.Debug(
		"label resized",
		slog.Int("width", message.Size.Width),
		slog.Int("height", message.Size.Height),
		slog.String("name", component.name),
	)
	component.size = message.Size
}
