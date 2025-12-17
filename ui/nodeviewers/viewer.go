package nodeviewers

import tea "github.com/charmbracelet/bubbletea"

type Viewer interface {
	Init() tea.Cmd
	UpdateViewer(tea.Msg) (Viewer, tea.Cmd)
	View() string
}

func UpdateViewerChild(child *Viewer, message tea.Msg, commands *[]tea.Cmd) {
	updatedChild, command := (*child).UpdateViewer(message)
	*child = updatedChild
	*commands = append(*commands, command)
}
