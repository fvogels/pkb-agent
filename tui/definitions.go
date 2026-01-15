package tui

import (
	"github.com/gdamore/tcell/v3"
)

const (
	SafeMode = true
)

type Component interface {
	Handle(Message)
	Render() Grid
}

type ComponentBase struct {
	Identifier   int
	Name         string
	MessageQueue MessageQueue
	Size         Size
}

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

type Position struct {
	X int
	Y int
}

type Size struct {
	Width  int
	Height int
}
