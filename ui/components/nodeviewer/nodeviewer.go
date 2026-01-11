package nodeviewer

import (
	"log/slog"
	"pkb-agent/pkg"
	"pkb-agent/ui/components/linksview"
	"pkb-agent/ui/components/nullviewer"
	"pkb-agent/ui/debug"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	size                           util.Size                   // Size of the component
	viewer                         tea.Model                   // Viewer specialized in node type
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
		util.UpdateUntypedChild(&model.viewer, message, &commands)
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

	util.UpdateUntypedChild(&model.viewer, tea.WindowSizeMsg{
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
		Links:     NewSliceAdapter(util.Map(node.GetLinks(), func(n *pkg.Node) string { return n.GetName() })),
		Backlinks: NewSliceAdapter(util.Map(node.GetBacklinks(), func(n *pkg.Node) string { return n.GetName() })),
	}, &commands)

	// Get viewer
	model.viewer = nil // message.Node.GetViewer()

	commands = append(commands, model.viewer.Init())

	updatedModel, extraCommands := model.updateChildSizes()
	model = updatedModel
	commands = append(commands, extraCommands)

	return model, tea.Batch(commands...)
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
