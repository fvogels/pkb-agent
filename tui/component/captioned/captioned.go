package captioned

import (
	"fmt"
	"pkb-agent/tui"
	"pkb-agent/tui/data"
	"pkb-agent/util/uid"
	"reflect"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	caption      data.Value[[]rune]
	captionStyle *tui.Style
	child        tui.Component
}

type MeasurableComponent struct {
	Component
}

func New(messageQueue tui.MessageQueue, caption data.Value[[]rune], child tui.Component) *Component {
	var component Component

	initialize(&component, messageQueue, caption, child)

	return &component
}

func initialize(component *Component, messageQueue tui.MessageQueue, caption data.Value[[]rune], child tui.Component) {
	captionStyle := tcell.StyleDefault.Background(color.Blue).Foreground(color.Reset)

	component.Identifier = uid.Generate()
	component.Name = "unnamed captioned"
	component.MessageQueue = messageQueue
	component.caption = caption
	component.captionStyle = &captionStyle
	component.child = child
}

func NewMeasurable(messageQueue tui.MessageQueue, caption data.Value[[]rune], child tui.MeasurableComponent) *MeasurableComponent {
	component := New(messageQueue, caption, child)

	return &MeasurableComponent{
		Component: *component,
	}
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

func (component *MeasurableComponent) Measure() tui.Size {
	child := component.child
	measurableChild, ok := child.(tui.MeasurableComponent)
	if !ok {
		panic(fmt.Sprintf("child with type %s is not measurable", reflect.TypeOf(component.child).String()))
	}

	childSize := measurableChild.Measure()

	return tui.Size{
		Width:  childSize.Width,
		Height: childSize.Height + 1,
	}
}
