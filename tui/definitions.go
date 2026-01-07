package tui

import "github.com/gdamore/tcell/v3"

const (
	SafeMode = true
)

type Component interface {
	Handle(Message)
	Render() Grid
}

type Message any

type Grid interface {
	GetSize() Size
	Get(position Position) Cell
}

type Style = tcell.Style

type Cell struct {
	Contents rune
	Style    *Style
	OnClick  func()
}

type MsgResize struct {
	Size Size
}

type MsgKey struct {
	Key string
}

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}
