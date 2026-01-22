package tui

import (
	"pkb-agent/tui/position"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
)

const (
	SafeMode = true
)

type Component interface {
	GetIdentifier() int
	Handle(Message)
	Render() Grid
}

type MeasurableComponent interface {
	Component
	Measure() size.Size
}

type ComponentBase struct {
	Identifier   int
	Name         string
	MessageQueue MessageQueue
	Size         size.Size
}

func (base *ComponentBase) GetIdentifier() int {
	return base.Identifier
}

type Grid interface {
	Size() size.Size
	At(position position.Position) Cell
}

type Style = tcell.Style

type Cell struct {
	Contents rune
	Style    *Style
	OnClick  func()
}
