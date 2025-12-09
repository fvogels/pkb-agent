package mainscreen

import tea "github.com/charmbracelet/bubbletea"

type mode interface {
	onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd)
}
