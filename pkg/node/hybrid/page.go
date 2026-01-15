package hybrid

import (
	"pkb-agent/pkg/node"
	"pkb-agent/tui"
)

type Page interface {
	GetCaption() string
	CreateViewer(tui.MessageQueue) tui.Component
	GetActions() []node.Action
}
