package hybrid

import tea "github.com/charmbracelet/bubbletea"

type msgSetSubviewers struct {
	recipient  int
	subviewers []tea.Model
}
