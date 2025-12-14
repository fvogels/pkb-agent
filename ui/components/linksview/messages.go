package linksview

import tea "github.com/charmbracelet/bubbletea"

type MsgSetLinks struct {
	Links     List
	Backlinks List
}

type msgLinksListWrapper struct {
	wrapped tea.Msg
}

type msgBacklinksListWrapper struct {
	wrapped tea.Msg
}
