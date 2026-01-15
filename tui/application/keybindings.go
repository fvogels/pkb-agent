package application

import "pkb-agent/tui"

var (
	BindingQuit = tui.KeyBinding{
		Key:         "q",
		Description: "quit",
		Message:     MsgQuit{},
	}

	BindingSelect = tui.KeyBinding{
		Key:         "Enter",
		Description: "select",
		Message:     MsgSelectHighlightedNode{},
	}

	BindingUnselect = tui.KeyBinding{
		Key:         "Delete",
		Description: "pop",
		Message:     MsgUnselectLastNode{},
	}

	BindingSearch = tui.KeyBinding{
		Key:         "/",
		Description: "search",
		Message:     MsgActivateInputMode{},
	}
)
