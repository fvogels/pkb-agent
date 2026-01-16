package backblaze

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
	"pkb-agent/tui/component/label"
	"pkb-agent/tui/data"
	"pkb-agent/ui/uid"
)

type Component struct {
	tui.ComponentBase
	child tui.Component
}

func NewViewer(messageQueue tui.MessageQueue) *Component {
	component := Component{
		ComponentBase: tui.ComponentBase{
			Identifier:   uid.Generate(),
			Name:         "unnamed backblaze viewer",
			MessageQueue: messageQueue,
		},
		child: label.New(messageQueue, "backblaze label", data.NewConstant("backblaze!")),
	}

	return &component
}

func (component *Component) Handle(message tui.Message) {
	switch message := message.(type) {
	case tui.MsgActivate:
		message.Respond(
			component.Identifier,
			component.onActivate,
			component.child,
		)

	case tui.MsgResize:
		component.onResize(message)
	}
}

func (component *Component) onActivate() {
	component.MessageQueue.Enqueue(messages.MsgSetNodeKeyBindings{
		Bindings: list.New[tui.KeyBinding](),
	})
}

func (component *Component) Render() tui.Grid {
	return component.child.Render()
}

func (component *Component) onResize(message tui.MsgResize) {
	component.Size = message.Size
	component.child.Handle(message)
}

// var keyMap = struct {
// 	Download key.Binding
// }{
// 	Download: key.NewBinding(
// 		key.WithKeys("d"),
// 		key.WithHelp("d", "download"),
// 	),
// }

// type Model struct {
// 	size util.Size
// 	node *RawNode
// }

// func NewViewer(node *RawNode) Model {
// 	return Model{
// 		node: node,
// 	}
// }

// func (model Model) Init() tea.Cmd {
// 	return model.signalKeybindingsUpdate()
// }

// func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
// 	return model.TypedUpdate(message)
// }

// func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
// 	switch message := message.(type) {
// 	case tea.WindowSizeMsg:
// 		return model.onResize(message)

// 	case tea.KeyMsg:
// 		return model.onKeyPressed(message)

// 	default:
// 		return model, nil
// 	}
// }

// func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
// 	model.size = util.Size{
// 		Width:  message.Width,
// 		Height: message.Height,
// 	}

// 	return model, nil
// }

// func (model Model) View() string {
// 	return fmt.Sprintf("%s/%s", model.node.bucket, model.node.filename)
// }

// func (model Model) signalDownloadAndOpen() tea.Cmd {
// 	return func() tea.Msg {
// 		if err := model.downloadAndOpen(); err != nil {
// 			panic("failed to download and open backblaze file")
// 		}
// 		return nil
// 	}
// }

// func (model Model) downloadAndOpen() error {
// 	bucketName := model.node.bucket
// 	filename := model.node.filename
// 	application_key := os.Getenv("APPLICATION_KEY")
// 	application_key_id := os.Getenv("APPLICATION_KEY_ID")
// 	ctx := context.Background()

// 	client, err := bb.New(ctx, application_key, application_key_id)
// 	if err != nil {
// 		slog.Error("Failed to create backblaze client")
// 		return err
// 	}

// 	buffer, err := client.DownloadToBuffer(ctx, bucketName, filename, 1, func(bytesDownloaded int) {})
// 	if err != nil {
// 		slog.Error(
// 			"Failed to download file from Backblaze",
// 			slog.String("bucket", bucketName),
// 			slog.String("filename", filename),
// 		)
// 		return err
// 	}

// 	zippedFiles, err := zipfile.Unpack(buffer)
// 	if err != nil {
// 		slog.Error("Failed to unzip files")
// 		return err
// 	}

// 	zippedFile := zippedFiles[0]
// 	path, err := zippedFile.SaveToDirectory(pathlib.New("."))
// 	if err != nil {
// 		slog.Error("Failed to save downloaded file")
// 		return err
// 	}

// 	if err := extern.OpenUsingDefaultViewer(path); err != nil {
// 		slog.Error("Failed to open downloaded file using default viewer")
// 		return err
// 	}

// 	return nil
// }

// func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
// 	switch {
// 	case key.Matches(message, keyMap.Download):
// 		return model.onDownloadAndOpen()

// 	default:
// 		return model, nil
// 	}
// }

// func (model Model) onDownloadAndOpen() (Model, tea.Cmd) {
// 	return model, model.signalDownloadAndOpen()
// }

// func (model Model) signalKeybindingsUpdate() tea.Cmd {
// 	return func() tea.Msg {
// 		return node.MsgUpdateNodeViewerBindings{
// 			KeyBindings: []key.Binding{
// 				keyMap.Download,
// 			},
// 		}
// 	}
// }
