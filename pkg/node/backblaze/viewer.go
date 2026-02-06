package backblaze

import (
	"fmt"
	"log/slog"
	"pkb-agent/extern"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/tui/grid"
	"pkb-agent/util/uid"
)

type Component struct {
	tui.ComponentBase
	root            tui.Component
	rawNode         *RawNode
	bindingDownload tui.KeyBinding
}

func NewViewer(messageQueue tui.MessageQueue, rawNode *RawNode) *Component {
	slog.Debug("Creating new backblaze viewer")

	identifier := uid.Generate()

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   identifier,
			Name:         "unnamed backblaze viewer",
			MessageQueue: messageQueue,
		},
		rawNode: rawNode,
		root:    createRoot(messageQueue),
		bindingDownload: tui.KeyBinding{
			Key:         "d",
			Description: "download",
			Message: tui.MsgCommand{
				Command: func() {
					go extern.BackblazeDownloadAndOpen(
						rawNode.bucket,
						rawNode.filename,
					)
				},
			},
		},
	}

	return &component
}

func createRoot(messageQueue tui.MessageQueue) tui.Component {
	return label.New(messageQueue, "backblaze label", data.NewConstant("backblaze!"))
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgStateUpdated:
		component.root.Handle(message)
		component.onStateUpdated()

	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgKey:
		component.onKey(message)

	default:
		component.root.Handle(message)
	}
}

func (component *Component) onKey(message tui.MsgKey) {
	tui.HandleKeyBindings(
		component.MessageQueue,
		message,
		component.bindingDownload,
	)
}

func (component *Component) onStateUpdated() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.FromItems(component.bindingDownload),
	})
}

func (component *Component) Render() grid.FiniteGrid {
	slog.Debug(
		"Rendering backblaze viewer",
		slog.String("size", component.Size.String()),
		slog.String("address", fmt.Sprintf("%p", component)),
	)

	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	slog.Debug(
		"Resizing backblaze viewer",
		slog.String("size", message.Size.String()),
		slog.String("address", fmt.Sprintf("%p", component)),
	)

	component.Size = message.Size
	component.root.Handle(message)
}
