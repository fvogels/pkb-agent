package hybrid

import (
	"fmt"
	"pkb-agent/graph/node"
	"pkb-agent/ui/uid"
	"pkb-agent/util"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	id                         int
	size                       util.Size
	node                       *RawNode
	data                       *nodeData // (strong) pointer to the node data, keeps information alive while viewer exists
	activePageIndex            int
	subviewers                 []tea.Model
	statusBarPageLocationStyle lipgloss.Style
	statusBarPageCaptionStyle  lipgloss.Style
	actionKeyBindings          []ActionKeyBinding
	pageActionKeyBindings      [][]ActionKeyBinding
}

type ActionKeyBinding struct {
	action     node.Action
	keyBinding key.Binding
}

var keyMap = struct {
	PreviousPage key.Binding
	NextPage     key.Binding
}{
	PreviousPage: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("s-tab", "previous"),
	),
	NextPage: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next"),
	),
}

func NewViewer(node *RawNode, data *nodeData) Model {
	actionKeyBindings, pageActionKeyBindings := createActionKeyBindings(data.actions, data.pages)

	return Model{
		id:                         uid.Generate(),
		node:                       node,
		data:                       data,
		statusBarPageLocationStyle: lipgloss.NewStyle().Background(lipgloss.Color("#88FF88")),
		statusBarPageCaptionStyle:  lipgloss.NewStyle().Background(lipgloss.Color("#AAFFAA")),
		actionKeyBindings:          actionKeyBindings,
		pageActionKeyBindings:      pageActionKeyBindings,
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

	return model, model.signalKeyBindingsUpdate()
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
	// Deal with node key bindings
	for _, actionKeyBinding := range model.actionKeyBindings {
		if key.Matches(message, actionKeyBinding.keyBinding) {
			return model, model.signalPerformAction(actionKeyBinding.action)
		}
	}

	// Deal with page key bindings
	if len(model.data.pages) > 0 {
		for _, actionKeyBinding := range model.pageActionKeyBindings[model.activePageIndex] {
			if key.Matches(message, actionKeyBinding.keyBinding) {
				return model, model.signalPerformAction(actionKeyBinding.action)
			}
		}
	}

	// Deal with fixed key bindings
	switch {
	case key.Matches(message, keyMap.PreviousPage):
		return model.onSwitchToPreviousPage()

	case key.Matches(message, keyMap.NextPage):
		return model.onSwitchToNextPage()

	default:
		return model, nil
	}
}

func (model Model) signalPerformAction(action node.Action) tea.Cmd {
	return func() tea.Msg {
		action.Perform()

		return nil
	}
}

func (model Model) View() string {
	if len(model.subviewers) != 0 {
		activeSubviewer := model.subviewers[model.activePageIndex]

		return lipgloss.JoinVertical(
			0,
			lipgloss.NewStyle().Height(model.size.Height-1).Render(activeSubviewer.View()),
			model.renderStatusBar(),
		)
	} else {
		return lipgloss.JoinVertical(
			0,
			lipgloss.NewStyle().Height(model.size.Height-1).Render(""),
			model.renderStatusBar(),
		)
	}
}

func (model Model) signalKeyBindingsUpdate() tea.Cmd {
	return func() tea.Msg {
		keyBindings := model.determineKeyBindings()

		return node.MsgUpdateNodeViewerBindings{
			KeyBindings: keyBindings,
		}
	}
}

func (model Model) determineKeyBindings() []key.Binding {
	keyBindings := []key.Binding{}

	// Add "previous page" and "next page" keybindings if there are at least two pages
	if len(model.subviewers) >= 2 {
		keyBindings = append(keyBindings, keyMap.PreviousPage, keyMap.NextPage)
	}

	// Add node action bindings
	for _, actionKeyBinding := range model.actionKeyBindings {
		keyBindings = append(keyBindings, actionKeyBinding.keyBinding)
	}

	// Add page action bindings
	if len(model.data.pages) > 0 {
		pageKeyBindings := model.pageActionKeyBindings[model.activePageIndex]

		for _, pageKeyBinding := range pageKeyBindings {
			keyBindings = append(keyBindings, pageKeyBinding.keyBinding)
		}
	}

	return keyBindings
}

func (model Model) renderStatusBar() string {
	if len(model.data.pages) > 0 {
		currentPage := model.activePageIndex + 1
		totalPageCount := len(model.data.pages)
		pageLocation := model.statusBarPageLocationStyle.Render(fmt.Sprintf(" Page %d/%d ", currentPage, totalPageCount))
		pageLocationWidth := lipgloss.Width(pageLocation)
		pageCaption := model.statusBarPageCaptionStyle.Width(model.size.Width - pageLocationWidth).Render(" " + model.data.pages[model.activePageIndex].GetCaption())

		return lipgloss.JoinHorizontal(0, pageLocation, pageCaption)
	} else {
		return model.statusBarPageLocationStyle.Width(model.size.Width).Render(" no pages ")
	}
}

func (model Model) onSwitchToPreviousPage() (Model, tea.Cmd) {
	if len(model.subviewers) > 1 {
		model.activePageIndex = (model.activePageIndex - 1 + len(model.subviewers)) % len(model.subviewers)
	}

	return model, nil
}

func (model Model) onSwitchToNextPage() (Model, tea.Cmd) {
	if len(model.subviewers) > 1 {
		model.activePageIndex = (model.activePageIndex + 1) % len(model.subviewers)
	}

	return model, nil
}

func createActionKeyBindings(actions []node.Action, pages []Page) ([]ActionKeyBinding, [][]ActionKeyBinding) {
	keys := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}
	nodeActionBindings := []ActionKeyBinding{}

	for index, action := range actions {
		keyBinding := key.NewBinding(
			key.WithKeys(keys[index]),
			key.WithHelp(keys[index], action.GetDescription()),
		)

		actionKeyBinding := ActionKeyBinding{
			keyBinding: keyBinding,
			action:     action,
		}

		nodeActionBindings = append(nodeActionBindings, actionKeyBinding)
	}

	pageActionBindings := [][]ActionKeyBinding{}
	startIndex := len(actions)

	for _, page := range pages {
		actionBindings := []ActionKeyBinding{}

		for actionIndex, action := range page.GetActions() {
			keyBinding := key.NewBinding(
				key.WithKeys(keys[startIndex+actionIndex]),
				key.WithHelp(keys[startIndex+actionIndex], action.GetDescription()),
			)

			actionKeyBinding := ActionKeyBinding{
				keyBinding: keyBinding,
				action:     action,
			}

			actionBindings = append(actionBindings, actionKeyBinding)
		}

		pageActionBindings = append(pageActionBindings, actionBindings)
	}

	return nodeActionBindings, pageActionBindings
}
