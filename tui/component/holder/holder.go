package holder

import (
	"log/slog"
	"pkb-agent/tui"
	"pkb-agent/tui/data"
)

type Component struct {
	size  tui.Size
	child data.Value[tui.Component]
}

func New(child data.Value[tui.Component]) *Component {
	component := Component{
		child: child,
	}

	// Make sure that whenever a new component is put in, it is resized
	child.Observe(func() {
		slog.Debug("!!!", "w", component.size.Width, "h", component.size.Height)
		c := child.Get()

		if c != nil {
			c.Handle(tui.MsgResize{
				Size: component.size,
			})
		}
	})

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) Render() tui.Grid {
	child := component.child.Get()

	if child != nil {
		return child.Render()
	} else {
		return tui.NewEmptyGrid(component.size)
	}
}

func (component *Component) onResize(message tui.MsgResize) {
	component.size = message.Size

	child := component.child.Get()
	if child != nil {
		child.Handle(message)
	}
}
