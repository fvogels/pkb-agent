package backblaze

import (
	"context"
	"log/slog"
	"os"
	bb "pkb-agent/backblaze"
	"pkb-agent/extern"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/util/pathlib"
	"pkb-agent/util/uid"
	"pkb-agent/util/zipfile"
)

type Component struct {
	tui.ComponentBase
	root            tui.Component
	rawNode         *RawNode
	bindingDownload tui.KeyBinding
}

func NewViewer(messageQueue tui.MessageQueue, rawNode *RawNode) *Component {
	identifier := uid.Generate()

	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   identifier,
			Name:         "unnamed backblaze viewer",
			MessageQueue: messageQueue,
		},
		rawNode: rawNode,
		root:    label.New(messageQueue, "backblaze label", data.NewConstant("backblaze!")),
		bindingDownload: tui.KeyBinding{
			Key:         "d",
			Description: "download",
			Message: tui.MsgCommand{
				Command: func() {
					go downloadAndOpen(
						rawNode.bucket,
						rawNode.filename,
					)
				},
			},
		},
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgActivate:
		message.Respond(
			component.Identifier,
			component.onActivate,
			component.root,
		)

	case tui.MsgResize:
		component.onResize(message)

	case tui.MsgKey:
		component.onKey(message)
	}
}

func (component *Component) onKey(message tui.MsgKey) {
	tui.HandleKeyBindings(
		component.MessageQueue,
		message,
		component.bindingDownload,
	)
}

func (component *Component) onActivate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.FromItems(component.bindingDownload),
	})
}

func (component *Component) Render() tui.Grid {
	return component.root.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.root.Handle(message)
}

func downloadAndOpen(bucketName string, filename string) error {
	application_key := os.Getenv("APPLICATION_KEY")
	application_key_id := os.Getenv("APPLICATION_KEY_ID")
	ctx := context.Background()

	client, err := bb.New(ctx, application_key, application_key_id)
	if err != nil {
		slog.Error("Failed to create backblaze client")
		return err
	}

	buffer, err := client.DownloadToBuffer(ctx, bucketName, filename, 1, func(bytesDownloaded int) {})
	if err != nil {
		slog.Error(
			"Failed to download file from Backblaze",
			slog.String("bucket", bucketName),
			slog.String("filename", filename),
		)
		return err
	}

	zippedFiles, err := zipfile.Unpack(buffer)
	if err != nil {
		slog.Error("Failed to unzip files")
		return err
	}

	zippedFile := zippedFiles[0]
	path, err := zippedFile.SaveToDirectory(pathlib.New("."))
	if err != nil {
		slog.Error("Failed to save downloaded file")
		return err
	}

	if err := extern.OpenUsingDefaultViewer(path); err != nil {
		slog.Error("Failed to open downloaded file using default viewer")
		return err
	}

	return nil
}
