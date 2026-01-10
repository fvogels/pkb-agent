package docknorth

import (
	"pkb-agent/tui"
)

type Component struct {
	size              tui.Size
	mainChild         tui.Component
	dockedChild       tui.Component
	dockedChildHeight int
}

func New(dockedChild tui.Component, mainChild tui.Component, dockedChildHeight int) *Component {
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
		boundary:        component.dockedChildHeight,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size

	dockedChildSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  message.Size.Width,
			Height: component.dockedChildHeight,
		},
	}
	component.dockedChild.Handle(dockedChildSizeMessage)

	mainChildSizeMessage := tui.MsgResize{
		Size: tui.Size{
			Width:  message.Size.Width,
			Height: component.size.Height - component.dockedChildHeight,
		},
	}
	component.mainChild.Handle(mainChildSizeMessage)
}
