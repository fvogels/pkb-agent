package page

import (
	"fmt"
	"pkb-agent/persistent/list"
	"pkb-agent/pkg/node"
	"pkb-agent/tui"
)

type Page interface {
	GetCaption() string
	CreateViewer(tui.MessageQueue) tui.Component
	GetActions() []node.Action
}

type MsgSetPageKeyBindings struct {
	Bindings list.List[tui.KeyBinding]
}

func (message MsgSetPageKeyBindings) String() string {
	return fmt.Sprintf("MsgSetPageKeyBindings[%s]", list.String(message.Bindings, func(b tui.KeyBinding) string { return b.Key }))
}
