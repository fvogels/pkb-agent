package hybrid

import tea "github.com/charmbracelet/bubbletea"

type MsgSetSubviewers struct {
	Recipient  int
	Subviewers []tea.Model
}
