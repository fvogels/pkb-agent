package tui

type Message any

type MsgActivate struct{}

type MsgResize struct {
	Size Size
}

type MsgKey struct {
	Key string
}

type MsgUpdateLayout struct{}
