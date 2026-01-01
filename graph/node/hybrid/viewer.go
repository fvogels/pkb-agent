package hybrid

import (
	"fmt"
	"pkb-agent/ui/uid"
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	id                   int
	size                 util.Size
	node                 *RawNode
	data                 *nodeData
	activeSubviewerIndex int
	subviewers           []tea.Model
}

func NewViewer(node *RawNode, data *nodeData) Model {
	return Model{
		id:   uid.Generate(),
		node: node,
		data: data,
	}
}

func (model Model) Init() tea.Cmd {
	commands := []tea.Cmd{}

	// Create subviewer for each page
	subviewers := make([]tea.Model, len(model.data.pages))
	for pageIndex, page := range model.data.pages {
		// Create subviewer
		viewer := page.CreateViewer()

		// Store subviewer
		subviewers[pageIndex] = viewer

		// Initialize subviewer
		commands = append(commands, viewer.Init())
	}

	// We cannot update the model here (since Init receives a copy), so we send ourselves a message
	commands = append(commands, model.signalUpdateSubviewers(subviewers))

	// Inform higher up component of updated key bindings
	commands = append(commands, model.signalKeybindingsUpdate())

	return tea.Batch(commands...)
}

func (model Model) signalUpdateSubviewers(subviewers []tea.Model) tea.Cmd {
	return func() tea.Msg {
		return msgSetSubviewers{
			recipient:  model.id,
			subviewers: subviewers,
		}
	}
}

func (model Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	return model.TypedUpdate(message)
}

func (model Model) TypedUpdate(message tea.Msg) (Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		return model.onResize(message)

	case tea.KeyMsg:
		return model.onKeyPressed(message)

	case msgSetSubviewers:
		if model.id == message.recipient {
			return model.onSetSubviewers(message)
		} else {
			return model, nil
		}

	default:
		commands := []tea.Cmd{}

		for subviewerIndex := range model.subviewers {
			util.UpdateUntypedChild(&model.subviewers[subviewerIndex], message, &commands)
		}

		return model, tea.Batch(commands...)
	}
}

func (model Model) onSetSubviewers(message msgSetSubviewers) (Model, tea.Cmd) {
	model.subviewers = message.subviewers

	return model, nil
}

func (model Model) onResize(message tea.WindowSizeMsg) (Model, tea.Cmd) {
	model.size = util.Size{
		Width:  message.Width,
		Height: message.Height,
	}

	commands := []tea.Cmd{}
	subviewerMessage := tea.WindowSizeMsg{
		Width:  model.size.Width,
		Height: model.size.Height - 1,
	}
	for subviewerIndex := range len(model.subviewers) {
		util.UpdateUntypedChild(&model.subviewers[subviewerIndex], subviewerMessage, &commands)
	}

	return model, tea.Batch(commands...)
}

func (model Model) onKeyPressed(message tea.KeyMsg) (Model, tea.Cmd) {
	switch {
	default:
		return model, nil
	}
}

func (model Model) View() string {
	if len(model.subviewers) != 0 {
		activeSubviewer := model.subviewers[model.activeSubviewerIndex]

		return lipgloss.JoinVertical(
			0,
			lipgloss.NewStyle().Height(model.size.Height-1).Render(activeSubviewer.View()),
			model.renderStatusBar(),
		)
	} else {
		return "no pages"
	}
}

func (model Model) signalKeybindingsUpdate() tea.Cmd {
	return func() tea.Msg {
		return nil
	}
}

func (model Model) renderStatusBar() string {
	currentPage := model.activeSubviewerIndex + 1
	totalPageCount := len(model.data.pages)
	return fmt.Sprintf("Page %d/%d", currentPage, totalPageCount)
}
