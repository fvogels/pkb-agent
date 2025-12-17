package nodeviewer

import (
	"log/slog"
	"pkb-agent/graph/nodes/atom"
	"pkb-agent/graph/nodes/backblaze"
	"pkb-agent/graph/nodes/bookmark"
	"pkb-agent/graph/nodes/markdown"
	"pkb-agent/graph/nodes/snippet"
	"pkb-agent/ui/components/linksview"
	"pkb-agent/ui/debug"
	"pkb-agent/ui/nodeviewers"
	"pkb-agent/ui/nodeviewers/bbviewer"
	"pkb-agent/ui/nodeviewers/bookmarkviewer"
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
	size      util.Size          // Size of the component
	viewer    nodeviewers.Viewer // Viewer specialized in node type
	linksView linksview.Model    // Links and backlinks view
}

func New() Model {
	return Model{
		viewer:    nullviewer.New(),
		linksView: linksview.New(),
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

func (model Model) UpdateViewer(message tea.Msg) (nodeviewers.Viewer, tea.Cmd) {
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

	commands := []tea.Cmd{}

	linksViewHeight := model.determineLinksViewHeight()

	util.UpdateChild(&model.linksView, tea.WindowSizeMsg{
		Width:  message.Width,
		Height: linksViewHeight,
	}, &commands)

	nodeviewers.UpdateViewerChild(&model.viewer, tea.WindowSizeMsg{
		Width:  message.Width,
		Height: message.Height - linksViewHeight - 1, // -1 needed for border
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

	// Get desired height and resize
	util.UpdateChild(&model.linksView, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.determineLinksViewHeight(),
	}, &commands)

	// Select correct viewer appropriate for node type
	switch nodeData := node.Extra.(type) {
	case *atom.Extra:
		model.viewer = nullviewer.New()

	case *snippet.Extra:
		model.viewer = snippetviewer.New(nodeData)

	case *bookmark.Extra:
		model.viewer = bookmarkviewer.New(nodeData)

	case *backblaze.Extra:
		model.viewer = bbviewer.New(nodeData)

	case *markdown.Extra:
		model.viewer = mdviewer.New(nodeData)

	default:
		slog.Debug(
			"unrecognized node type",
			slog.String("type", reflect.TypeOf(node.Extra).String()),
		)

		model.viewer = nullviewer.New()
	}

	commands = append(commands, model.viewer.Init())
	nodeviewers.UpdateViewerChild(&model.viewer, tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.size.Height,
	}, &commands)

	return model, tea.Batch(commands...)
}

// determineLinksViewHeight tries to give the links view its desired height,
// but still keeps some room for the node viewer
func (model *Model) determineLinksViewHeight() int {
	desiredHeight := model.linksView.GetDesiredHeight()
	return util.MinInt(model.size.Height-11, desiredHeight)
}

func (model *Model) GetKeyBindings() []key.Binding {
	result := []key.Binding{}

	return result
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
