package application

import (
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
)

type MsgQuit struct{}

type MsgSelectHighlightedNode struct{}

type MsgSetModeKeyBindings struct {
	Bindings list.List[tui.KeyBinding]
}

type MsgSetNodeKeyBindings struct {
	Bindings list.List[tui.KeyBinding]
}

type MsgActivateInputMode struct{}

type MsgActivateViewMode struct{}

type MsgActivateMode struct{}

type MsgUnselectLastNode struct{}
