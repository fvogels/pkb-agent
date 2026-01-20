package messages

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/tui"
)

type MsgQuit struct{}

func (message MsgQuit) String() string {
	return "MsgQuit"
}

type MsgSelectHighlightedNode struct{}

func (message MsgSelectHighlightedNode) String() string {
	return "MsgSelectHighlightedNode"
}

type MsgSetModeKeyBindings struct {
	Bindings list.List[tui.KeyBinding]
}

func (message MsgSetModeKeyBindings) String() string {
	return fmt.Sprintf("MsgSetModeKeyBindings[Bindings=%s]", list.String(message.Bindings, func(b tui.KeyBinding) string { return b.Key }))
}

type MsgSetNodeKeyBindings struct {
	Bindings list.List[tui.KeyBinding]
}

func (message MsgSetNodeKeyBindings) String() string {
	return fmt.Sprintf("MsgSetNodeKeyBindings[Bindings=%s]", list.String(message.Bindings, func(b tui.KeyBinding) string { return b.Key }))
}

type MsgActivateInputMode struct{}

func (message MsgActivateInputMode) String() string {
	return "MsgActivateInputMode"
}

type MsgActivateViewMode struct{}

func (message MsgActivateViewMode) String() string {
	return "MsgActivateViewMode"
}

type MsgUnselectLastNode struct{}

func (message MsgUnselectLastNode) String() string {
	return "MsgUnselectLastNode"
}

type MsgSwitchLinksView struct{}

func (message MsgSwitchLinksView) String() string {
	return "MsgSwitchLinksView"
}

type MsgLockSelectedNodes struct{}

func (message MsgLockSelectedNodes) String() string {
	return "MsgLockSelectedNodes"
}

type MsgUnlockSelectedNodes struct{}

func (message MsgUnlockSelectedNodes) String() string {
	return "MsgUnlockSelectedNodes"
}
