package bbviewer

import (
	"context"
	"log/slog"
	"os"
	bb "pkb-agent/backblaze"
	"pkb-agent/extern"
	"pkb-agent/graph/nodes/backblaze"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"pkb-agent/zipfile"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var keyMap = struct {
	Download key.Binding
}{
	Download: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "download"),
	),
}

type Model struct {
	size                           util.Size
	nodeData                       *backblaze.Extra
	status                         status
	createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg
}

func New(createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg, nodeData *backblaze.Extra) Model {
	return Model{
		nodeData:                       nodeData,
		createUpdateKeyBindingsMessage: createUpdateKeyBindingsMessage,
	}
}

func (model Model) Init() tea.Cmd {
	return func() tea.Msg {
		return model.createUpdateKeyBindingsMessage([]key.Binding{
			keyMap.Download,
		})
	}
}

func (model Model) UpdateViewer(message tea.Msg) (nodeviewers.Viewer, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

	case msgStatusUpdate:
		model.status = message.status
		return model, model.signalListen(model.status.getChannel())
	}

	return model, nil
}

func (model Model) View() string {
	parts := []string{}
	parts = append(parts, model.nodeData.Filename)

	if model.status != nil {
		parts = append(parts, model.status.view())
	}

	return lipgloss.JoinVertical(0, parts...)
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model, nil
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, keyMap.Download):
		return model.openFile()

	default:
		return model, nil
	}
}

func (model Model) openFile() (Model, tea.Cmd) {
	channel := make(chan status)

	go func() {
		defer close(channel)

		application_key := os.Getenv("APPLICATION_KEY")
		application_key_id := os.Getenv("APPLICATION_KEY_ID")
		ctx := context.Background()

		client, err := bb.New(ctx, application_key, application_key_id)
		if err != nil {
			slog.Error("Failed to create backblaze client")
			channel <- statusErrorOccurred{err: err}
			return
		}

		buffer, err := client.DownloadToBuffer(ctx, model.nodeData.BucketName, model.nodeData.Filename, 1, func(bytesDownloaded int) {
			channel <- statusDownloading{
				bytesDownloaded: bytesDownloaded,
				channel:         channel,
			}
		})
		if err != nil {
			slog.Error(
				"Failed to download file from Backblaze",
				slog.String("bucket", model.nodeData.BucketName),
				slog.String("filename", model.nodeData.Filename),
			)
			channel <- statusErrorOccurred{err: err}
			return
		}

		channel <- statusUnzipping{
			channel: channel,
		}
		zippedFiles, err := zipfile.Unpack(buffer)
		if err != nil {
			slog.Error("Failed to unzip files")
			channel <- statusErrorOccurred{err: err}
			return
		}

		zippedFile := zippedFiles[0]
		path, err := zippedFile.SaveToDirectory(pathlib.New("."))
		if err != nil {
			slog.Error("Failed to save downloaded file")
			channel <- statusErrorOccurred{err: err}
			return
		}

		if err := extern.OpenUsingDefaultViewer(path); err != nil {
			slog.Error("Failed to open downloaded file using default viewer")
			channel <- statusErrorOccurred{err: err}
			return
		}

		channel <- statusFinished{}
	}()

	return model, model.signalListen(channel)
}

func (model Model) signalListen(channel chan status) tea.Cmd {
	if channel != nil {
		return func() tea.Msg {
			status := <-channel
			return msgStatusUpdate{
				status: status,
			}
		}
	} else {
		return nil
	}
}
