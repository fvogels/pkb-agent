package mainscreen

import (
	"pkb-agent/util"

	tea "github.com/charmbracelet/bubbletea"
)

type mode interface {
	onKeyPressed(model Model, message tea.KeyMsg) (Model, tea.Cmd)

	activate(model *Model) tea.Cmd
	resize(model *Model, size util.Size) tea.Cmd
	render(model *Model) string
}
