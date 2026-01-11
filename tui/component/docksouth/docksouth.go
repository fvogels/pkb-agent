package docksouth

import (
	"pkb-agent/tui"
)

type Component struct {
	size              tui.Size
	mainChild         tui.Component
	dockedChild       tui.Component
	dockedChildHeight int
}

func New(mainChild tui.Component, dockedChild tui.Component, dockedChildHeight int) *Component {
	return &Component{
		mainChild:         mainChild,
		dockedChild:       dockedChild,
		dockedChildHeight: dockedChildHeight,
	}
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.mainChild.Handle(message)
		component.dockedChild.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	return &grid{
		size:            component.size,
		mainChildGrid:   component.mainChild.Render(),
		dockedChildGrid: component.dockedChild.Render(),
		boundary:        component.size.Height - component.dockedChildHeight,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size
	component.updateLayout()
}

func (component *Component) updateLayout() {
	width := component.size.Width
	dockedChildHeight := component.dockedChildHeight
	mainChildHeight := component.size.Height - component.dockedChildHeight

	dockedChildSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  width,
			Height: dockedChildHeight,
		},
	}
	component.dockedChild.Handle(dockedChildSizeMessage)

	mainChildSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  width,
			Height: mainChildHeight,
		},
	}
	component.mainChild.Handle(mainChildSizeMessage)
}
