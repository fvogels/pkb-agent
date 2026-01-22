package vstack

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/size"
	"pkb-agent/util/uid"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Component struct {
	tui.ComponentBase
	children   list.List[tui.MeasurableComponent]
	emptyStyle *tui.Style
}

func New(messageQueue tui.MessageQueue, children list.List[tui.MeasurableComponent]) *Component {
	emptyStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed vstack",
			MessageQueue: messageQueue,
		},
		children:   children,
		emptyStyle: &emptyStyle,
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		for index := range component.children.Size() {
			child := component.children.At(index)
			child.Handle(message)
		}
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size

	for index := range component.children.Size() {
		child := component.children.At(index)
		childRequestedSize := child.Measure()

		child.Handle(tui.MsgResize{
			Size: size.Size{
				Width:  message.Size.Width,
				Height: childRequestedSize.Height,
			},
		})
	}
}

func (component *Component) Render() tui.Grid {
	childGrids := list.MapList(
		component.children,
		func(child tui.MeasurableComponent) tui.Grid {
			return child.Render()
		},
	)

	return newGrid(component, childGrids)
}
