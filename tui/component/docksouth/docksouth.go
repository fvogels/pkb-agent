package docksouth

import (
	"log/slog"
	"pkb-agent/tui"
	"pkb-agent/tui/debug"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	mainChild         tui.Component
	dockedChild       tui.Component
	dockedChildHeight int
}

func New(messageQueue tui.MessageQueue, name string, mainChild tui.Component, dockedChild tui.Component, dockedChildHeight int) *Component {
	return &Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         name,
			MessageQueue: messageQueue,
		},
		mainChild:         mainChild,
		dockedChild:       dockedChild,
		dockedChildHeight: dockedChildHeight,
	}
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)

	default:
		component.mainChild.Handle(message)
		component.dockedChild.Handle(message)
	}
}

func (component *Component) Render() tui.Grid {
	slog.Debug("Rendering docksouth", slog.String("name", component.Name))

	return &grid{
		size:            component.Size,
		mainChildGrid:   component.mainChild.Render(),
		dockedChildGrid: component.dockedChild.Render(),
		boundary:        component.Size.Height - component.dockedChildHeight,
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.updateLayout()
}

func (component *Component) updateLayout() {
	width := component.Size.Width
	dockedChildHeight := component.dockedChildHeight
	mainChildHeight := component.Size.Height - component.dockedChildHeight

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
