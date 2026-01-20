package captioned

import (
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	caption      data.Value[[]rune]
	captionStyle *tui.Style
	child        tui.Component
}

func New(messageQueue *tui.MessageQueue, caption data.Value[[]rune], child tui.Component) *Component {
	captionStyle := tcell.StyleDefault.Background(color.Blue).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed captioned",
			MessageQueue: *messageQueue,
		},
		caption:      caption,
		captionStyle: &captionStyle,
		child:        child,
	}

	return &component
}

func (component *Component) SetCaptionStyle(style *tui.Style) {
	component.captionStyle = style
}

func (component *Component) Render() tui.Grid {
	return newGrid(component)
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.child.Handle(message)
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	childSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  message.Size.Width,
			Height: component.Size.Height - 1,
		},
	}
	component.child.Handle(childSizeMessage)
}
