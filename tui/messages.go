package tui

import (
	"pkb-agent/persistent/list"
)

type Message any

type MsgResize struct {
	Size Size
}

type MsgKey struct {
	Key string
}

type MsgUpdateLayout struct{}

type MsgSetNodeKeyBindings struct {
	Bindings list.List[KeyBinding]
}
