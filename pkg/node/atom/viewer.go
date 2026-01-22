package atom

import (
	"fmt"
	"log/slog"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/tui/debug"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	child tui.Component
}

func NewViewer(messageQueue tui.MessageQueue) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed atom viewer",
			MessageQueue: messageQueue,
		},
		child: label.New(messageQueue, "atom label", data.NewConstant("atom!")),
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	debug.LogMessage(message)

	switch message := message.(type) {
	case tui.MsgStateUpdated:
		component.child.Handle(message)
		component.onStateUpdated()

	case tui.MsgResize:
		component.onResize(message)

	default:
		component.child.Handle(message)
	}
}

func (component *Component) onStateUpdated() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})
}

func (component *Component) Render() grid.Grid {
	slog.Debug(
		"Rendering atom view",
		slog.String("size", component.Size.String()),
		slog.String("address", fmt.Sprintf("%p", component)),
	)

	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	slog.Debug(
		"Resizing atom viewer",
		slog.String("size", message.Size.String()),
		slog.String("address", fmt.Sprintf("%p", component)),
	)

	component.Size = message.Size
	component.child.Handle(message)
}
