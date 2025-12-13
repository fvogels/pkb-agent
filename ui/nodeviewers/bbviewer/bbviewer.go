package bbviewer

import (
	"context"
	"log/slog"
	"os"
	bb "pkb-agent/backblaze"
	"pkb-agent/extern"
	"pkb-agent/graph/nodes/backblaze"
	"pkb-agent/ui/debug"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"pkb-agent/zipfile"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var keyMap = struct {
	OpenInBrowser key.Binding
}{
	OpenInBrowser: key.NewBinding(
		key.WithKeys("o"),
		key.WithHelp("o", "open"),
	),
}

type Model struct {
	size     util.Size
	nodeData *backblaze.Extra
}

func New(nodeData *backblaze.Extra) Model {
	return Model{
		nodeData: nodeData,
	}
}

func (model Model) Init() tea.Cmd {
	return nil
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)
	}

	return model, nil
}

func (model Model) View() string {
	return model.nodeData.Filename
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
	case key.Matches(message, keyMap.OpenInBrowser):
		return model.openFile()

	default:
		return model, nil
	}
}

func (model Model) openFile() (Model, tea.Cmd) {
	command := func() tea.Msg {
		application_key := os.Getenv("APPLICATION_KEY")
		application_key_id := os.Getenv("APPLICATION_KEY_ID")
		ctx := context.Background()

		client, err := bb.New(ctx, application_key, application_key_id)
		if err != nil {
			slog.Error("Failed to create backblaze client")
			panic("failed to create backblaze client")
		}

		buffer, err := client.DownloadToBuffer(ctx, model.nodeData.BucketName, model.nodeData.Filename, 1, func(n int) {})
		if err != nil {
			slog.Error(
				"Failed to download file from BackBlaze",
				slog.String("bucket", model.nodeData.BucketName),
				slog.String("filename", model.nodeData.Filename),
			)
			panic("failed to download file")
		}

		zippedFiles, err := zipfile.Unpack(buffer)
		if err != nil {
			slog.Error("Failed to unzip files")
			panic("failed to unzip files")
		}

		zippedFile := zippedFiles[0]
		path, err := zippedFile.SaveToDirectory(pathlib.New("."))
		if err != nil {
			return err
		}

		if err := extern.OpenUsingDefaultViewer(path); err != nil {
			return err
		}

		return nil
	}

	return model, command
}
