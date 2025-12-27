package hybridviewer

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"pkb-agent/backblaze"
	"pkb-agent/extern"
	"pkb-agent/graph/nodes/hybrid"
	"pkb-agent/ui/components/markdownview"
	"pkb-agent/ui/components/sourceviewer"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/util"
	"pkb-agent/util/pathlib"
	"pkb-agent/zipfile"
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size                           util.Size
	nodeInfo                       *hybrid.Info
	nodeData                       *hybrid.Data
	pageViewers                    []tea.Model
	activePage                     int
	commands                       []command
	createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg
	statusBarStyle                 lipgloss.Style
}

var keyMap = struct {
	SwitchToNextPage key.Binding
}{
	SwitchToNextPage: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next page"),
	),
}

type PageViewer interface {
	PageViewerUpdate(tea.Msg) (PageViewer, tea.Cmd)
	View() string
}

type command struct {
	keyBinding key.Binding
	perform    func(*Model) tea.Cmd
}

func New(createUpdateKeyBindingsMessage func(keyBindings []key.Binding) tea.Msg, nodeData *hybrid.Info) Model {
	return Model{
		nodeInfo:                       nodeData,
		pageViewers:                    nil,
		activePage:                     -1,
		createUpdateKeyBindingsMessage: createUpdateKeyBindingsMessage,
		statusBarStyle:                 lipgloss.NewStyle().Background(lipgloss.Color("#AAAAFF")),
	}
}

func (model Model) Init() tea.Cmd {
	return model.signalLoadNodeData()
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) UpdateViewer(message tea.Msg) (nodeviewers.Viewer, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case msgMarkdownLoaded:
		return model.onDataLoaded(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

	default:
		commands := []tea.Cmd{}

		for pageViewerIndex := range model.pageViewers {
			util.UpdateUntypedChild(&model.pageViewers[pageViewerIndex], message, &commands)
		}

		return model, tea.Batch(commands...)
	}
}

func (model Model) View() string {
	var viewerResult string
	if model.activePage == -1 {
		viewerResult = ""
	} else {
		height := model.size.Height - 1
		viewerResult = lipgloss.NewStyle().Height(height).MaxHeight(height).Render(model.pageViewers[model.activePage].View())
	}

	statusBar := model.renderStatusBar()

	slog.Debug("rendered hybrid node", "viewerHeight", lipgloss.Height(viewerResult), "statusBarHeight", lipgloss.Height(statusBar), "totalHeight", model.size.Height)

	return lipgloss.JoinVertical(0, viewerResult, statusBar)
}

func (model Model) renderStatusBar() string {
	activePage := model.activePage
	totalPages := len(model.pageViewers)

	var contents string
	if totalPages > 0 {
		contents = fmt.Sprintf("Page %d/%d", activePage+1, totalPages)
	} else {
		contents = "no pages"
	}

	return model.statusBarStyle.Width(model.size.Width).Render(contents)
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	commands := []tea.Cmd{}

	// Size of each viewer
	// Decrease height by one to allow for status bar
	viewerSize := util.Size{
		Width:  message.Width,
		Height: message.Height - 1,
	}

	for index := range model.pageViewers {
		util.UpdateUntypedChild(&model.pageViewers[index], viewerSize, &commands)
	}

	return model, tea.Batch(commands...)
}

func (model *Model) signalLoadNodeData() tea.Cmd {
	info := model.nodeInfo

	return func() tea.Msg {
		data, err := info.GetData()
		if err != nil {
			slog.Error("Error while reading node data", slog.String("error", err.Error()))
			panic("failed to load node's data")
		}

		return msgMarkdownLoaded{
			data,
		}
	}
}

func (model Model) onDataLoaded(message msgMarkdownLoaded) (Model, tea.Cmd) {
	model.nodeData = message.data
	model.commands = model.createCommands()
	commands := []tea.Cmd{model.signalUpdatedKeyBindings()}
	pageViewers := []tea.Model{}

	for _, page := range model.nodeData.Pages {
		model.activePage = 0

		switch page := page.(type) {
		case *hybrid.MarkdownPage:
			viewer := markdownview.New()
			command1 := viewer.Init()

			var command2 tea.Cmd
			viewer, command2 = viewer.TypedUpdate(markdownview.MsgSetSource{
				Source: page.Source,
			})

			commands = append(commands, command1, command2)
			pageViewers = append(pageViewers, viewer)

		case *hybrid.SnippetPage:
			viewer := sourceviewer.New()
			command1 := viewer.Init()

			var command2 tea.Cmd
			viewer, command2 = viewer.TypedUpdate(sourceviewer.MsgSetSource{
				Source:   page.Source,
				Language: page.Language,
			})

			commands = append(commands, command1, command2)
			pageViewers = append(pageViewers, viewer)

		default:
			slog.Error("Unknown page type", slog.String("pageType", reflect.TypeOf(page).String()))
			panic("unknown page type")
		}
	}

	model.pageViewers = pageViewers

	return model, tea.Batch(commands...)
}

func (model Model) signalUpdatedKeyBindings() tea.Cmd {
	return func() tea.Msg {
		return model.createUpdateKeyBindingsMessage(model.determineKeyBindings())
	}
}

func (model Model) createCommands() []command {
	commands := []command{}
	keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	for actionIndex, action := range model.nodeData.Actions {
		command := createCommand(keys[actionIndex], action)
		commands = append(commands, command)
	}

	return commands
}

func (model Model) determineKeyBindings() []key.Binding {
	return util.Map(model.commands, func(action command) key.Binding {
		return action.keyBinding
	})
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	case key.Matches(message, keyMap.SwitchToNextPage):
		return model.onSwitchToNextPage()

	default:
		for _, action := range model.commands {
			if key.Matches(message, action.keyBinding) {
				return model, action.perform(&model)
			}
		}

		return model, nil
	}
}

func (model Model) onSwitchToNextPage() (Model, tea.Cmd) {
	if model.activePage == -1 {
		// No pages available; do nothing
		return model, nil
	}

	model.activePage = (model.activePage + 1) % len(model.pageViewers)

	return model, nil
}

func createCommand(boundKey string, action hybrid.Action) command {
	binding := key.NewBinding(
		key.WithKeys(boundKey),
		key.WithHelp(boundKey, action.GetDescription()),
	)

	switch action := action.(type) {
	case *hybrid.BrowserAction:
		return createBrowserCommand(binding, action)

	case *hybrid.DownloadAction:
		return createDownloadCommand(binding, action)

	default:
		slog.Error("Unrecognized action", slog.String("actionType", reflect.TypeOf(action).String()))
		panic("unrecognized action")
	}
}

func createBrowserCommand(binding key.Binding, action *hybrid.BrowserAction) command {
	return command{
		keyBinding: binding,
		perform: func(model *Model) tea.Cmd {
			return func() tea.Msg {
				openURL(action.URL)
				return nil
			}
		},
	}
}

func createDownloadCommand(binding key.Binding, action *hybrid.DownloadAction) command {
	return command{
		keyBinding: binding,
		perform: func(model *Model) tea.Cmd {
			return func() tea.Msg {
				downloadAndOpenBackblazeFile(action.Bucket, action.Filename)
				return nil
			}
		},
	}
}

func openURL(url string) {
	if err := extern.OpenURLInBrowser(url); err != nil {
		slog.Error("Failed to open URL in browser", slog.String("url", url))
		panic("failed to open browser")
	}
}

func downloadAndOpenBackblazeFile(bucket string, filename string) {
	application_key := os.Getenv("APPLICATION_KEY")
	application_key_id := os.Getenv("APPLICATION_KEY_ID")
	ctx := context.Background()

	client, err := backblaze.New(ctx, application_key, application_key_id)
	if err != nil {
		slog.Error("Failed to create backblaze client")
		return
	}

	buffer, err := client.DownloadToBuffer(ctx, bucket, filename, 1, func(bytesDownloaded int) {})
	if err != nil {
		slog.Error(
			"Failed to download file from Backblaze",
			slog.String("bucket", bucket),
			slog.String("filename", filename),
		)
		return
	}

	zippedFiles, err := zipfile.Unpack(buffer)
	if err != nil {
		slog.Error("Failed to unzip files")
		return
	}

	zippedFile := zippedFiles[0]
	path, err := zippedFile.SaveToDirectory(pathlib.New("."))
	if err != nil {
		slog.Error("Failed to save downloaded file")
		return
	}

	if err := extern.OpenUsingDefaultViewer(path); err != nil {
		slog.Error("Failed to open downloaded file using default viewer")
		return
	}
}
