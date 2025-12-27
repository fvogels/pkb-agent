package nodeviewer

import (
	"log/slog"
	"pkb-agent/graph/nodes/atom"
	"pkb-agent/graph/nodes/backblaze"
	"pkb-agent/graph/nodes/bookmark"
	"pkb-agent/graph/nodes/hybrid"
	"pkb-agent/graph/nodes/markdown"
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/ui/components/linksview"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/ui/nodeviewers/bbviewer"
	"pkb-agent/ui/nodeviewers/bookmarkviewer"
	"pkb-agent/ui/nodeviewers/hybridviewer"
	"pkb-agent/ui/nodeviewers/mdviewer"
	"pkb-agent/ui/nodeviewers/nullviewer"
	"pkb-agent/ui/nodeviewers/snippetviewer"
	"pkb-agent/util"
	"reflect"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size                           util.Size                   // Size of the component
	viewer                         nodeviewers.Viewer          // Viewer specialized in node type
	linksView                      linksview.Model             // Links and backlinks view
	createUpdateKeyBindingsMessage func([]key.Binding) tea.Msg // Used by viewers to inform the main screen that key bindings need an update
}

func New(createUpdateKeyBindingsMessage func([]key.Binding) tea.Msg) Model {
	return Model{
		viewer:                         nullviewer.New(),
		linksView:                      linksview.New(),
		createUpdateKeyBindingsMessage: createUpdateKeyBindingsMessage,
	}
}

func (model Model) Init() tea.Cmd {
	return tea.Batch(
		model.viewer.Init(),
		model.linksView.Init(),
	)
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	debug.ShowBubbleTeaMessage(message)

	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResized(message)

	case MsgSetNode:
		return model.onSetNode(message)

	default:
		commands := []tea.Cmd{}
		nodeviewers.UpdateViewerChild(&model.viewer, message, &commands)
		util.UpdateChild(&model.linksView, message, &commands)
		return model, tea.Batch(commands...)
	}
}

func (model Model) View() string {
	return lipgloss.JoinVertical(
		0,
		lipgloss.NewStyle().Border(lipgloss.ASCIIBorder(), false, false, true, false).Render(
			model.linksView.View(),
		),
		model.viewer.View(),
	)
}

func (model Model) onResized(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	return model.updateChildSizes()
}

// updateChildSizes is called when the layout of the child components have to be updated.
func (model Model) updateChildSizes() (Model, tea.Cmd) {
	parentWidth := model.size.Width
	parentHeight := model.size.Height
	commands := []tea.Cmd{}

	linksViewHeight := model.determineLinksViewHeight()
	specializedViewerHeight := parentHeight - linksViewHeight - 1 // -1 needed for border separating links and node view

	slog.Debug(
		"Resizing nodeviewer",
		slog.Int("totalHeight", model.size.Height),
		slog.Int("linksViewerHeight", linksViewHeight),
		slog.Int("viewerHeight", specializedViewerHeight),
	)

	util.UpdateChild(&model.linksView, tea.WindowSizeMsg{
		Width:  parentWidth,
		Height: linksViewHeight,
	}, &commands)

	nodeviewers.UpdateViewerChild(&model.viewer, tea.WindowSizeMsg{
		Width:  parentWidth,
		Height: specializedViewerHeight,
	}, &commands)

	return model, tea.Batch(commands...)
}

func (model Model) onSetNode(message MsgSetNode) (Model, tea.Cmd) {
	node := message.Node
	commands := []tea.Cmd{}

	// Update links and backlinks
	util.UpdateChild(&model.linksView, linksview.MsgSetLinks{
		Links:     NewSliceAdapter(node.Links),
		Backlinks: NewSliceAdapter(node.Backlinks),
	}, &commands)

	// Select correct viewer appropriate for node type
	switch nodeData := node.Info.(type) {
	case *atom.Info:
		model.viewer = nullviewer.New()

	case *snippet.Info:
		model.viewer = snippetviewer.New(model.createUpdateKeyBindingsMessage, nodeData)

	case *bookmark.Info:
		model.viewer = bookmarkviewer.New(model.createUpdateKeyBindingsMessage, nodeData)

	case *backblaze.Info:
		model.viewer = bbviewer.New(model.createUpdateKeyBindingsMessage, nodeData)

	case *markdown.Info:
		model.viewer = mdviewer.New(nodeData)

	case *hybrid.Info:
		model.viewer = hybridviewer.New(model.createUpdateKeyBindingsMessage, nodeData)

	default:
		slog.Debug(
			"unrecognized node type",
			slog.String("type", reflect.TypeOf(node.Info).String()),
		)

		model.viewer = nullviewer.New()
	}

	commands = append(commands, model.viewer.Init())

	updatedModel, extraCommands := model.updateChildSizes()
	commands = append(commands, extraCommands)

	return updatedModel, tea.Batch(commands...)
}

// determineLinksViewHeight tries to give the links view its desired height,
// but still keeps some room for the node viewer
func (model *Model) determineLinksViewHeight() int {
	desiredHeight := model.linksView.GetDesiredHeight()
	return util.MinInt(model.size.Height-11, desiredHeight)
}

type SliceAdapter[T any] struct {
	slice []T
}

func NewSliceAdapter[T any](slice []T) *SliceAdapter[T] {
	return &SliceAdapter[T]{
		slice: slice,
	}
}

func (adapter *SliceAdapter[T]) Length() int {
	return len(adapter.slice)
}

func (adapter *SliceAdapter[T]) At(index int) T {
	return adapter.slice[index]
}
