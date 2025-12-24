package helpbar

import "github.com/charmbracelet/bubbles/key"

type MsgSetKeyBindings struct {
	KeyBindings []key.Binding
}

type MsgSetNodeCounts struct {
	Remaining int
	Total     int
}
