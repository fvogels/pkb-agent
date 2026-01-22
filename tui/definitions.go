package tui

import (
	"pkb-agent/tui/grid"
	"pkb-agent/tui/size"

	"github.com/gdamore/tcell/v3"
)

const (
	SafeMode = true
)

type Component interface {
	GetIdentifier() int
	Handle(Message)
	Render() grid.Grid
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

type Style = tcell.Style
