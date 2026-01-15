package tui

type Message any

type MsgResize struct {
	Size Size
}

type MsgKey struct {
	Key string
}

type MsgUpdateLayout struct{}
