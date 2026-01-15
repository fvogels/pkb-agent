package docknorth

import (
	"log/slog"
	"pkb-agent/tui"
	"pkb-agent/ui/uid"
)

type Component struct {
	tui.ComponentBase
	mainChild         tui.Component
	dockedChild       tui.Component
	dockedChildHeight int
}

func New(messageQueue tui.MessageQueue, name string, dockedChild tui.Component, mainChild tui.Component, dockedChildHeight int) *Component {
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

func (component *Component) SetDockerChildHeight(height int) {
	component.dockedChildHeight = height
	component.updateLayout()
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
		size:            component.Size,
		mainChildGrid:   component.mainChild.Render(),
		dockedChildGrid: component.dockedChild.Render(),
		boundary:        component.dockedChildHeight,
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
			Height: component.dockedChildHeight,
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

	slog.Debug(
		"updated layout in docknorth",
		slog.String("name", component.Name),
		slog.Int("totalHeight", component.Size.Height),
		slog.Int("dockedChildHeight", dockedChildHeight),
		slog.Int("mainChildHeight", mainChildHeight),
	)
}
