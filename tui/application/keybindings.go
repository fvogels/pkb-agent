package application

import (
	"pkb-agent/tui"
	"pkb-agent/tui/application/messages"
)

var (
	BindingQuit = tui.KeyBinding{
		Key:         "q",
		Description: "quit",
		Message:     messages.MsgQuit{},
	}

	BindingSelect = tui.KeyBinding{
		Key:         "Enter",
		Description: "select",
		Message:     messages.MsgSelectHighlightedNode{},
	}

	BindingUnselect = tui.KeyBinding{
		Key:         "Delete",
		Description: "pop",
		Message:     messages.MsgUnselectLastNode{},
	}

	BindingSearch = tui.KeyBinding{
		Key:         "/",
		Description: "search",
		Message:     messages.MsgActivateInputMode{},
	}

	BindingSwitchLinksView = tui.KeyBinding{
		Key:         "*",
		Description: "links",
		Message:     messages.MsgSwitchLinksView{},
	}
)
